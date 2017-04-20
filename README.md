# Gim

Gim is tool for Gopher teams 

## Installation

Download and install Gim: 
```
go get github.com/retargetapp/gim
```
Gim will be installed to `$GOPATH/bin` directory.

**Notice:** We recommend to add `$GOPATH/bin` to your `$PATH` environment variable and use command: `gim`, instead of command: `$GOPATH/bin/gim`   

## Set up Gim in your project

First of all you need to create a directory with migration sources files on your project directory.
Example:
```bash
mkdir migrations 
```

To set up your project Gim config and install Git's hooks run:

```
gim init
```

During `git init` execution you should set:
1. DB driver you use, example:
```
mysql
```
2. [DSN](https://en.wikipedia.org/wiki/Data_source_name) in follow format: `user:password@tcp(host:port)/db_name`, example: 
```
gim:EDS#TG$@tcp(127.0.0.1:3306)/gim
```
3. directory with migration sources files, example
```
./migrations
```

`gim init` creates `.gim.yml` file to store config and installs some hooks into `.git/hooks/post-merge` and `.git/hooks/post-checkout` files.

## Migrations sources format

Every migration version consists of two sql-scripts: `up` and `down` stored in separated files: `<version>.up.sql` and `<version>.down.sql`.
`<version>` is a timestamp of migration created moment. Both files can not be empty. 

**Notice:** We recommend use one version file per sql DDL query. This approach helps to resolve problems with incompatible migrations faster.
 

## Usage


**Notice:** You `$GOAPTH` enviropment variable should be defined correctly every time when you use Gim