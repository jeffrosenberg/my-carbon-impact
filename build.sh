#! /usr/bin/env sh
# TODO: Replace with Make and/or CodeBuild

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
GO_BIN="$SCRIPT_DIR/go/bin"
pushd $SCRIPT_DIR > /dev/null
pushd go/cmd > /dev/null

for entity in $(ls -1); do
  pushd "$entity" > /dev/null
  for operation in $(ls -1); do
    pushd "$operation" > /dev/null
      echo "building ${entity}-${operation}"
      GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o "${GO_BIN}/${entity}-${operation}"
    popd > /dev/null
  done
  popd > /dev/null
done

popd > /dev/null
popd > /dev/null