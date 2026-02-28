# Maintainer: bb8TheDev1717
pkgname=master-pkg
pkgver=0.2.0
pkgrel=1
pkgdesc="Unified package manager CLI for pacman and AUR"
arch=('x86_64')
url="https://github.com/bb8TheDev1717/Master-package-installer"
license=('MIT')
depends=('pacman' 'paru')
makedepends=('go')
source=("$pkgname-$pkgver.tar.gz::$url/archive/v$pkgver.tar.gz")
sha256sums=('a77faea52157539ddd32062d4d4ea30dd0aa55d1b50203acb3012f123073d7dd')

build() {
    cd "Master-package-installer-$pkgver"
    go build -o build/master .
}

package() {
    cd "Master-package-installer-$pkgver"
    install -Dm755 build/master "$pkgdir/usr/bin/master"
}
