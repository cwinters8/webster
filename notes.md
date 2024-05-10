# notes

Value proposition: gives freelancers a method to enable their clients to update static content without hand-holding.

Enable users to create and manage website content in a WYSIWYG editor, like Frontpage used to be, and provide previews and deployments.

## differentiators

Note: Many of the 'features' below are not in place yet.

### Data ownership

You own all your data by default. Your content files are accessible and you can do with them as you wish.

### Offline first

Your content files are stored on your local device and synced to the cloud in the background when a connection is available, so you can continue to work on your content even if you're in an area with spotty internet.

### Transparency

You can see and edit the underlying HTML if you want to. This also allows you to incrementally learn HTML, because you can see exactly how the page content changes as you make edits.

### Version control

Did something not quite look right on your recent deployment? You can roll any page back to a specific version.

### Integration

You can integrate pages not built with Webster in your website, enabling nearly endless possibilities 🚀

## Hosting

Host on Oracle Cloud [Container Instances](https://www.oracle.com/cloud/cloud-native/container-instances/) for free - each tenancy gets ~4 OCPUs and ~24 GB memory free every month for 'Ampere' (arm) compute

## Useful packages

[GitHub module for Go](https://github.com/google/go-github)
[git module](https://pkg.go.dev/github.com/go-git/go-git/v5) to enable any git provider

## tasks

- [ ] Pass existing content to editor
- [ ] Stay on the same page after POST request
- [ ] Authentication
- [ ] Save output HTML content to a GitHub repo and/or file store (e.g. S3 or an S3-compatible service)
- [ ] [Custom skin](https://www.tiny.cloud/docs/tinymce/latest/creating-a-skin/) for editor
