# Maintainer: bb8TheDev1717
pkgname=master-pkg
pkgver=0.1.0
pkgrel=1
pkgdesc="Unified package manager CLI for pacman and AUR"
arch=('x86_64')
url="https://github.com/bb8TheDev1717/Master-package-installer"
license=('MIT')
depends=('pacman' 'paru')
makedepends=('go')
source=("$pkgname-$pkgver.tar.gz::$url/archive/v$pkgver.tar.gz")
sha256sums=('SKIP')

build() {
    cd "master-$pkgver"
    go build -o master .
}

package() {
    cd "master-$pkgver"
    install -Dm755 master "$pkgdir/usr/bin/master"
}
