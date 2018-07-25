package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	humanize "github.com/dustin/go-humanize"
	minio "github.com/pierre-emmanuelJ/minio-go"

	"github.com/spf13/cobra"
)

const (
	printDate = "2006-01-02 15:04:05 MST"
)

// sosListCmd represents the list command
var sosListCmd = &cobra.Command{
	Use:     "list [<bucket name>/path]",
	Short:   "List file and folder",
	Aliases: gListAlias,
	RunE: func(cmd *cobra.Command, args []string) error {

		///XXX you must use a default zone support SOS
		minioClient, err := newMinioClient(gCurrentAccount.DefaultZone)
		if err != nil {
			println("XXX wainting for all zone supporting SOS: use a supported defaultZone")
			log.Fatal(err)
		}

		isRec, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			return err
		}

		if len(args) == 0 {
			return displayBucket(minioClient, isRec)
		}

		path := filepath.ToSlash(args[0])
		path = strings.Trim(path, "/")
		p := splitPath(args[0])

		if len(p) == 0 {
			return displayBucket(minioClient, isRec)
		}
		if p[0] == "" {
			return displayBucket(minioClient, isRec)
		}

		var prefix string
		if len(p) > 1 {
			prefix = path[len(p[0]):]
			prefix = strings.Trim(prefix, "/")
		}

		bucketName := p[0]

		///XXX waiting for pithos 301 redirect
		location, err := minioClient.GetBucketLocation(bucketName)
		if err != nil {
			return err
		}

		minioClient, err = newMinioClient(location)
		if err != nil {
			return err
		}
		///

		table := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight)

		if isRec {
			listRecursively(minioClient, bucketName, prefix, "", false, table)
			return table.Flush()
		}

		doneCh := make(chan struct{})
		defer close(doneCh)
		recursive := true

		for message := range minioClient.ListObjectsV2(bucketName, prefix, recursive, doneCh) {
			sPrefix := splitPath(prefix)
			sKey := splitPath(message.Key)
			if sPrefix[0] == "" && len(sKey) > 1 {
				continue
			}
			if len(sKey) > len(sPrefix)+1 {
				continue
			}
			if isPrefix(prefix, message.Key) {
				continue
			}
			lastModified := message.LastModified.Format(printDate)
			key := filepath.ToSlash(message.Key)
			key = strings.TrimLeft(key[len(prefix):], "/")

			fmt.Fprintf(table, "%s\t%s\t%s\n", // nolint: errcheck
				fmt.Sprintf("[%s]", lastModified),
				fmt.Sprintf("%6s ", humanize.IBytes(uint64(message.Size))),
				key) // nolint: errcheck
		}

		return table.Flush()
	},
}

func listRecursively(c *minio.Client, bucketName, prefix, zone string, displayBucket bool, table io.Writer) {
	doneCh := make(chan struct{})
	defer close(doneCh)

	for message := range c.ListObjectsV2(bucketName, prefix, true, doneCh) {

		lastModified := message.LastModified.Format(printDate)
		if displayBucket {
			fmt.Fprintf(table, "%s\t%s\t%s\t%s\n", fmt.Sprintf("[%s]", lastModified), // nolint: errcheck
				fmt.Sprintf("[%s]", zone),
				fmt.Sprintf("%6s ", humanize.IBytes(uint64(message.Size))),
				fmt.Sprintf("%s/%s", bucketName, message.Key)) // nolint: errcheck
		} else {
			fmt.Fprintf(table, "%s\t%s\t%s\n", fmt.Sprintf("[%s]", lastModified), // nolint: errcheck
				fmt.Sprintf("%6s ", humanize.IBytes(uint64(message.Size))),
				fmt.Sprintf("%s/%s", bucketName, message.Key)) // nolint: errcheck
		}
	}
}

func displayBucket(minioClient *minio.Client, isRecursive bool) error {
	allBuckets, err := listBucket(minioClient)
	if err != nil {
		return err
	}

	table := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight)

	for zoneName, buckets := range allBuckets {
		for _, bucket := range buckets {
			if isRecursive {
				///XXX Waiting for pithos 301 redirect
				minioClient, err = newMinioClient(zoneName)
				if err != nil {
					return err
				}
				///
				listRecursively(minioClient, bucket.Name, "", zoneName, true, table)
				continue
			}
			fmt.Fprintf(table, "%s\t%s\t%s\t%s/\n", // nolint: errcheck
				fmt.Sprintf("[%s]", bucket.CreationDate.Format(printDate)),
				fmt.Sprintf("[%s]", zoneName),
				fmt.Sprintf("%6s ", humanize.IBytes(uint64(0))),
				bucket.Name) // nolint: errcheck
		}
	}
	return table.Flush()
}

func listBucket(minioClient *minio.Client) (map[string][]minio.BucketInfo, error) {
	bucketInfos, err := minioClient.ListBuckets()
	if err != nil {
		return nil, err
	}

	res := map[string][]minio.BucketInfo{}

	for _, bucketInfo := range bucketInfos {

		bucketLocation, err := minioClient.GetBucketLocation(bucketInfo.Name)
		if err != nil {
			return nil, err
		}
		if _, ok := res[bucketLocation]; !ok {
			res[bucketLocation] = []minio.BucketInfo{bucketInfo}
			continue
		}

		res[bucketLocation] = append(res[bucketLocation], bucketInfo)

	}
	return res, nil
}

func isPrefix(prefix, file string) bool {
	prefix = strings.Trim(prefix, "/")
	file = strings.Trim(file, "/")
	return prefix == file
}

func splitPath(s string) []string {
	path := filepath.ToSlash(s)
	path = strings.Trim(path, "/")
	return strings.Split(path, "/")
}

func init() {
	sosCmd.AddCommand(sosListCmd)
	sosListCmd.Flags().BoolP("recursive", "r", false, "List recursively")
}
