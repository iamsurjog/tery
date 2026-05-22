# Maintainer: Sujatro Ganguli <iamsurjog@gmail.com>
pkgname=tery
pkgver=1.0.1
pkgrel=1
pkgdesc="Lightweight battery monitoring daemon for Linux with configurable notifications"
arch=('any')
url="https://github.com/iamsurjog/tery"
license=('MIT')
optdepends=('libnotify')
makedepends=('go')

source=("git+${url}.git")
sha256sums=('SKIP')

build() {
    cd "${pkgname}"
    go build -o "${pkgname}" .
}

package() {
  cd "${pkgname}"
  # Create the destination directory first
  install -Dm755 tery "$pkgdir/usr/bin/tery"
}
