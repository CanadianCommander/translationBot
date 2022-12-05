#!/bin/bash

pushd $(dirname $0)/../

helm upgrade --install tb ./helm/.

popd