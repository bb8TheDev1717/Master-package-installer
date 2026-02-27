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
sha256sums=('94da4df3f865819e967c4926b423d48a9bb0eeaee0b16885352c14c6186f6f7f')

build() {
    cd "Master-package-installer-$pkgver"
    go build -o master .
}

package() {
    cd "Master-package-installer-$pkgver"
    install -Dm755 master "$pkgdir/usr/bin/master"
}
