class Emm < Formula
  desc "Eidolon Minion Manager - Modular Go-based CLI/TUI for AI"
  homepage "https://github.com/meesfatels/EMM"
  url "https://github.com/meesfatels/EMM.git", branch: "main"
  version "1.0.0"
  head "https://github.com/meesfatels/EMM.git", branch: "dev"
  license "MIT"

  depends_on "go" => :build

  def install
    # Use version from the formula or build from source
    ver = build.head? ? "dev-head" : version
    system "go", "build", "-trimpath", "-ldflags", "-s -w -X main.version=#{ver}", "-o", bin/"emm", "./cmd/emm"
  end

  def caveats
    <<~EOS
      After installation, run 'emm init' to set up your configuration in ~/.emm/
      You'll need an OpenRouter API key set in ~/.emm/emm.yaml to use it.
    EOS
  end

  test do
    # Simple check to see if the binary executes
    assert_match "emm", shell_output("#{bin}/emm --help")
  end
end
