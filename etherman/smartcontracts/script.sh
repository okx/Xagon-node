#!/bin/sh

set -e

gen() {
    local package=$1

    abigen --bin bin/${package}.bin --abi abi/${package}.abi --pkg=${package} --out=${package}/${package}.go
}

gen pol
gen mockverifier
gen proxy
gen preetrogpolygonzkevmglobalexitroot
gen preetrogpolygonzkevmbridge
gen preetrogpolygonzkevm
gen etrogpolygonzkevm
gen etrogpolygonzkevmglobalexitroot
gen etrogpolygonrollupmanager
gen mocketrogpolygonrollupmanager
gen elderberrypolygonzkevm
gen feijoapolygonzkevm
gen feijoapolygonzkevmglobalexitroot
gen feijoapolygonrollupmanager
gen mockfeijoapolygonrollupmanager
