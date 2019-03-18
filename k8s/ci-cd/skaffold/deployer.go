package deployer

//package build
type Artifact struct {
	ImageName string
	Tag       string
}

type Deployer interface {
	Labels() map[string]string

	// Deploy should ensure that the build results are deployed to the Kubernetes
	// cluster.
	Deploy(context.Context, io.Writer, []build.Artifact) ([]Artifact, error)

	// Dependencies returns a list of files that the deployer depends on.
	// In dev mode, a redeploy will be triggered
	Dependencies() ([]string, error)

	// Cleanup deletes what was deployed by calling Deploy.
	Cleanup(context.Context, io.Writer) error
}

type KubectlDeployer struct {
	*latest.KubectlDeploy

	workingDir  string
	kubectl     kubectl.CLI
	defaultRepo string
}

func (k *KubectlDeployer) Deploy(ctx context.Context, out io.Writer, builds []build.Artifact) ([]Artifact, error) {
	if err := k.kubectl.CheckVersion(); err != nil {
		color.Default.Fprintln(out, err)
	}

	//get all the yaml files
	//it support read yaml from another/remote cluster
	manifests, err := k.readManifests(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "reading manifests")
	}

	if len(manifests) == 0 {
		return nil, nil
	}

	//replace image in yaml file using default repo
	manifests, err = manifests.ReplaceImages(builds, k.defaultRepo)
	if err != nil {
		return nil, errors.Wrap(err, "replacing images in manifests")
	}

	//use kubectl appy to create the resources
	updated, err := k.kubectl.Apply(ctx, out, manifests)
	if err != nil {
		return nil, errors.Wrap(err, "apply")
	}

	return parseManifestsForDeploys(k.kubectl.Namespace, updated)
}
