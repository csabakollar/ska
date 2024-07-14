package contentprovider

import (
	"fmt"
	"github.com/apex/log"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"os"
)

type GitHub struct {
	remoteURI     string
	repositoryURL string
	repositoryRef string
	workingDir    string
	log           *log.Entry
}

const GitHubPrefix = "https://github.com/"

func NewGitHub(remoteURI string) (*GitHub, error) {
	logCtx := log.WithFields(log.Fields{
		"pkg":  "contentprovider",
		"type": "github",
	})
	tmpDir, err := os.MkdirTemp("", workingDirPrefix)
	if err != nil {
		return nil, err
	}
	return &GitHub{remoteURI: remoteURI, workingDir: tmpDir, log: logCtx}, nil
}

func (cp *GitHub) WorkingDir() string {
	return cp.workingDir
}

func (cp *GitHub) DownloadContent() error {
	if err := cp.validateRemoteURI(cp.remoteURI); err != nil {
		return err
	}
	token := os.Getenv("GITHUB_TOKEN")
	var auth *http.BasicAuth
	if token != "" {
		auth = &http.BasicAuth{
			Username: "x-access-token",
			Password: token,
		}
	}
	// check if the reference is a tag ot branch
	var referenceName plumbing.ReferenceName
	isTag, err := cp.remoteTagExists(cp.repositoryRef, token)
	if err != nil {
		return err
	}

	if isTag {
		referenceName = plumbing.ReferenceName("refs/tags/" + cp.repositoryRef)
	} else {
		referenceName = plumbing.ReferenceName("refs/heads/" + cp.repositoryRef)
	}

	_, err = git.PlainClone(cp.workingDir, false, &git.CloneOptions{
		URL:           cp.repositoryURL,
		ReferenceName: referenceName,
		Auth:          auth,
	})

	if err != nil {
		switch err.Error() {
		case "authentication required":
			return fmt.Errorf("to clone this repository you need to set the GITHUB_TOKEN environment variable, with a valid GitHub Personal Access Token. Error: %w", err)
		default:
			return err
		}
	}

	// remove the .git folder
	if err := cp.removeGitFolder(); err != nil {
		return err
	}
	return nil
}

func (cp *GitHub) Cleanup() error {
	cp.log.WithFields(log.Fields{"workingDir": cp.workingDir}).Debug("removing working dir.")
	err := os.RemoveAll(cp.workingDir)
	return err
}

func (cp *GitHub) RemoteURI() string {
	return cp.remoteURI
}
