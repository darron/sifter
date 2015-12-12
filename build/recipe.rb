class Sifter < FPM::Cookery::Recipe
  name 'sifter'

  version '0.7'
  revision '1'
  description 'sifter'

  homepage 'https://github.com/darron/sifter'
  source "https://github.com/darron/sifter/releases/download/v#{version}/sifter-#{version}-linux-amd64.zip"
  sha256 '367c98a025b7ff344b4e6cff14652acbcaab37ff06f795ac85b9ddef4a0cce96'

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
