pkgname=ttymer
pkgver=1.0.2
pkgrel=1
pkgdesc="A simple terminal-based timer written in Go with TUI support"
arch=('x86_64')
url="https://github.com/darwincereska/ttymer"
license=("MIT")
depends=("libnotify")
makedepends=("go")

source=("$pkgname-$pkgver.tar.gz::$url/archive/v$pkgver.tar.gz")
sha256sums=("SKIP")

build() {
    cd "$srcdir/$pkgname-$pkgver"
    export CGO_CPPFLAGS="${CPPFLAGS}"
    export CGO_CFLAGS="${CFLAGS}"
    export CGO_CXXFLAGS="${CXXFLAGS}"
    export CGO_LDFLAGS="${LDFLAGS}"
    export GOFLAGS="-buildmode=pie -trimpath -ldflags=-linkmode=external -mod=readonly -modcacherw"
    go build -o ttymer
}

package() {
    cd "$srcdir/$pkgname-$pkgver"
    install -Dm755 ttymer "$pkgdir/usr/bin/ttymer"

    install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"

    install -Dm644 README.md "$pkgdir/usr/share/doc/$pkgname/README.md"
}
