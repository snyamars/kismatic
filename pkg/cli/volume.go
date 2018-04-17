package cli

import (
	"io"

	"github.com/spf13/cobra"
)

// NewCmdVolume returns the storage command
func NewCmdVolume(in io.Reader, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volume",
		Short: "manage storage volumes on your Kubernetes cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}
	cmd.AddCommand(NewCmdVolumeAdd(out))
	cmd.AddCommand(NewCmdVolumeList(out))
	cmd.AddCommand(NewCmdVolumeDelete(in, out))
	return cmd
}
