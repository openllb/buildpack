module github.com/openllb/buildpack

go 1.12

require (
	github.com/moby/buildkit v0.6.3
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.4.0 // indirect
	golang.org/x/sys v0.0.0-20191026070338-33540a1f6037 // indirect
)

replace github.com/moby/buildkit => github.com/hinshun/buildkit v0.0.0-20200123030914-aacaae031fb3

replace github.com/hashicorp/go-immutable-radix => github.com/tonistiigi/go-immutable-radix v0.0.0-20170803185627-826af9ccf0fe

replace github.com/jaguilar/vt100 => github.com/tonistiigi/vt100 v0.0.0-20190402012908-ad4c4a574305
