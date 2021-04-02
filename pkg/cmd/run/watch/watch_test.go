package watch

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/cli/cli/pkg/cmdutil"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
)

func TestNewCmdWatch(t *testing.T) {
	tests := []struct {
		name     string
		cli      string
		tty      bool
		wants    WatchOptions
		wantsErr bool
	}{
		{
			name:     "blank nontty",
			wantsErr: true,
		},
		{
			name: "blank tty",
			tty:  true,
			wants: WatchOptions{
				Prompt:   true,
				Interval: 2,
			},
		},
		{
			name: "interval",
			tty:  true,
			cli:  "-i10",
			wants: WatchOptions{
				Interval: 10,
				Prompt:   true,
			},
		},
		{
			name: "exit status",
			cli:  "1234 --exit-status",
			wants: WatchOptions{
				Interval:   2,
				RunID:      "1234",
				ExitStatus: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, _, _ := iostreams.Test()
			io.SetStdinTTY(tt.tty)
			io.SetStdoutTTY(tt.tty)

			f := &cmdutil.Factory{
				IOStreams: io,
			}

			argv, err := shlex.Split(tt.cli)
			assert.NoError(t, err)

			var gotOpts *WatchOptions
			cmd := NewCmdWatch(f, func(opts *WatchOptions) error {
				gotOpts = opts
				return nil
			})
			cmd.SetArgs(argv)
			cmd.SetIn(&bytes.Buffer{})
			cmd.SetOut(ioutil.Discard)
			cmd.SetErr(ioutil.Discard)

			_, err = cmd.ExecuteC()
			if tt.wantsErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			assert.Equal(t, tt.wants.RunID, gotOpts.RunID)
			assert.Equal(t, tt.wants.Prompt, gotOpts.Prompt)
			assert.Equal(t, tt.wants.ExitStatus, gotOpts.ExitStatus)
			assert.Equal(t, tt.wants.Interval, gotOpts.Interval)
		})
	}
}

// TODO execution tests
