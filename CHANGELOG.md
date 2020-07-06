# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## Unreleased

## [0.2.0] - 2020-07-06

- Create new dir `deploy/olm-catalog/postgresql-operator/manifests` which the latest version which is point out for the next release version 0.2.0
- Allow define the Storage Class of the PVC used for the database
- Start to use CRDs version v1 instead of v1beta1
- Remove the support for k8s clusters < 1.16.0
- Upgrade project to use version 0.18.1 of SDK
- Update the title of the project in the CVS file

## [0.1.1] - 2019-09-09
- Add status in CSV OLM integration files 

## [0.1.0] - 2019-09-08
- Initial development release which allows work with as standalone and has OLM files. 
