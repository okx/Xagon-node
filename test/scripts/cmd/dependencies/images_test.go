package dependencies

import (
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_image_readCurrentDigest(t *testing.T) {
	var appFs = afero.NewMemMapFs()

	tcs := []struct {
		description      string
		input            string
		expectedOutput   string
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			description: "single container with matching image and existing digest",

			input: `version: '3'
services:
    zkevm-node:
        image: imageorg/imagerepo@sha256:f7bc95017b64a6dee98dd2e3a98cbea8c715de137d0f599b1b16f683c2dae955
`,
			expectedOutput: "sha256:f7bc95017b64a6dee98dd2e3a98cbea8c715de137d0f599b1b16f683c2dae955",
		},
		{
			description: "multiple containers with matching image and existing digest",
			input: `version: '3'
services:
    zkevm-node:
        image: zkevm-node

    x1-mock-l1-network:
        image: imageorg/imagerepo@sha256:f7bc95017b64a6dee98dd2e3a98cbea8c715de137d0f599b1b16f683c2dae955
    zkevm-mock-prover:
        image: hezprover
`,
			expectedOutput: "sha256:f7bc95017b64a6dee98dd2e3a98cbea8c715de137d0f599b1b16f683c2dae955",
		},
		{
			description: "single container with matching image and non-existing digest",

			input: `version: '3'
services:
    zkevm-node:
        image: imageorg/imagerepo:latest
`,
			expectedError:    true,
			expectedErrorMsg: "image \"imageorg/imagerepo\" not found",
		},
		{
			description: "single container with non-matching image",

			input: `version: '3'
services:
    zkevm-node:
        image: imageNeworg/imageNewrepo@sha256:f7bc95017b64a6dee98dd2e3a98cbea8c715de137d0f599b1b16f683c2dae955
`,
			expectedError:    true,
			expectedErrorMsg: "image \"imageorg/imagerepo\" not found",
		},
		{
			description:      "invalid yaml",
			input:            "not valid yaml",
			expectedError:    true,
			expectedErrorMsg: "yaml: unmarshal errors",
		},
	}

	const (
		defaultPath = "/a/b/dockerCompose.yml"
		imageName   = "imageorg/imagerepo"
	)

	subject := &imageUpdater{
		fs:             appFs,
		targetFilePath: defaultPath,
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			require.NoError(t, createFile(appFs, defaultPath, tc.input))
			defer func() {
				require.NoError(t, appFs.Remove(defaultPath))
			}()

			actualOutput, err := subject.readCurrentDigest(imageName)
			require.NoError(t, testutils.CheckError(err, tc.expectedError, tc.expectedErrorMsg))
			require.Equal(t, tc.expectedOutput, actualOutput)
		})
	}
}

func Test_image_updateDigest(t *testing.T) {
	var appFs = afero.NewMemMapFs()

	tcs := []struct {
		description               string
		initialFileContents       string
		oldDigest                 string
		newDigest                 string
		expectedFinalFileContents string
		expectedError             bool
		expectedErrorMsg          string
	}{
		{
			description: "single container with matching image and existing digest",
			initialFileContents: `version: '3'
services:
    zkevm-node:
        image: imageorg/imagerepo@sha256:oldDigest`,
			oldDigest: "sha256:oldDigest",
			newDigest: "sha256:newDigest",
			expectedFinalFileContents: `version: '3'
services:
    zkevm-node:
        image: imageorg/imagerepo@sha256:newDigest`,
		},
		{
			description: "single container, not matching image",
			initialFileContents: `version: '3'
services:
    zkevm-node:
        image: imageorg/anotherimagerepo@sha256:oldDigest`,
			oldDigest: "sha256:oldDigest",
			newDigest: "sha256:newDigest",
			expectedFinalFileContents: `version: '3'
services:
    zkevm-node:
        image: imageorg/anotherimagerepo@sha256:oldDigest`,
		},
		{
			description: "single container with matching image, non-existing digest",
			initialFileContents: `version: '3'
services:
    zkevm-node:
        image: imageorg/imagerepo@sha256:veryOldDigest`,
			oldDigest: "sha256:oldDigest",
			newDigest: "sha256:newDigest",
			expectedFinalFileContents: `version: '3'
services:
    zkevm-node:
        image: imageorg/imagerepo@sha256:veryOldDigest`,
		},
		{
			description: "multiple container with matching image and existing digest",
			initialFileContents: `version: '3'
services:
    x1-mock-l1-network:
        image: imageorg/networkImagerepo@sha256:oldDigest
    zkevm-node:
        image: imageorg/imagerepo@sha256:oldDigest
    zkevm-mock-prover:
        image: imageorg/proverImagerepo@sha256:oldDigest
`,
			oldDigest: "sha256:oldDigest",
			newDigest: "sha256:newDigest",
			expectedFinalFileContents: `version: '3'
services:
    x1-mock-l1-network:
        image: imageorg/networkImagerepo@sha256:oldDigest
    zkevm-node:
        image: imageorg/imagerepo@sha256:newDigest
    zkevm-mock-prover:
        image: imageorg/proverImagerepo@sha256:oldDigest
`,
		},
		{
			description:               "invalid yaml",
			initialFileContents:       "not valid yaml",
			oldDigest:                 "sha256:oldDigest",
			newDigest:                 "sha256:newDigest",
			expectedFinalFileContents: "not valid yaml",
		},
	}

	const (
		defaultPath = "/a/b/dockerCompose.yml"
		imageName   = "imageorg/imagerepo"
	)

	subject := &imageUpdater{
		fs:             appFs,
		targetFilePath: defaultPath,
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			require.NoError(t, createFile(appFs, defaultPath, tc.initialFileContents))
			defer func() {
				require.NoError(t, appFs.Remove(defaultPath))
			}()

			err := subject.updateDigest(imageName, tc.oldDigest, tc.newDigest)
			require.NoError(t, testutils.CheckError(err, tc.expectedError, tc.expectedErrorMsg))
			actualFileContents, err := afero.ReadFile(appFs, defaultPath)
			require.NoError(t, err)
			require.Equal(t, tc.expectedFinalFileContents, string(actualFileContents))
		})
	}
}

func createFile(appFs afero.Fs, path, content string) error {
	f, err := appFs.Create(path)

	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Errorf("Could not close %q: %v", path, err)
		}
	}()

	_, err = f.WriteString(content)

	return err
}
