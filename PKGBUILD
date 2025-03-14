# Maintainer: Your Name <your.email@example.com>
pkgname=tuiodo
pkgver=1.0.0
pkgrel=1
pkgdesc="A modern terminal task manager with extensive customization"
arch=('x86_64')
url="https://github.com/spmfte/tuiodo"
license=('MIT')
depends=()
makedepends=('go')
source=("$pkgname-$pkgver.tar.gz::https://github.com/spmfte/tuiodo/archive/v$pkgver.tar.gz")
sha256sums=('SKIP')

build() {
  cd "$pkgname-$pkgver"
  export CGO_CPPFLAGS="${CPPFLAGS}"
  export CGO_CFLAGS="${CFLAGS}"
  export CGO_CXXFLAGS="${CXXFLAGS}"
  export CGO_LDFLAGS="${LDFLAGS}"
  export GOFLAGS="-buildmode=pie -trimpath -ldflags=-linkmode=external -mod=readonly -modcacherw"
  go build -o tuiodo
}

package() {
  cd "$pkgname-$pkgver"
  install -Dm755 tuiodo "$pkgdir/usr/bin/tuiodo"
  install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
  install -Dm644 README.md "$pkgdir/usr/share/doc/$pkgname/README.md"
} 