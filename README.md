# Capn Hook

Easily configure and run git hooks to check your code

## Install

    go get github.com/dcu/capn-hook

## Generate manifest

To generate the default manifest type:

	$ capn-hook generate

and then to review the manifest type:

	$ cat hooks.yml


## Install

To install the hooks type:

	$ capn-hook install


## Run

The hooks will automatically run when git hooks are triggered. You can also run the hooks manually by typing:

	$ capn-hook run <hook name>

For example:

	$ capn-hook run pre-commit
