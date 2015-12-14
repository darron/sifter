class Sifter < FPM::Cookery::Recipe
  name 'sifter'

  version '0.7'
  revision '5'
  description 'sifter'

  homepage 'https://github.com/darron/sifter'
  source "https://github.com/darron/sifter/releases/download/v#{version}/sifter-#{version}-linux-amd64.zip"
  sha256 'f56329cae73810fa81173f87ca4e002b7244b00308bbf534afac7d37ac0467ae'

  maintainer 'Darron <darron@froese.org>'
  vendor 'octohost'

  license 'Mozilla Public License, version 2.0'

  conflicts 'sifter'
  replaces 'sifter'

  build_depends 'unzip'

  def build
    safesystem "mkdir -p #{builddir}/usr/local/bin/"
    safesystem "cp -f #{builddir}/sifter-#{version}-linux-amd64/sifter-#{version}-linux-amd64 #{builddir}/usr/local/bin/sifter"
  end

  def install
    safesystem "mkdir -p #{destdir}/usr/local/bin/"
    safesystem "cp -f #{builddir}/usr/local/bin/sifter #{destdir}/usr/local/bin/sifter"
    safesystem "chmod 755 #{destdir}/usr/local/bin/sifter"
  end
end
