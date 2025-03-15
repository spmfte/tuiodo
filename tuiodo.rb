class Tuiodo < Formula
  desc "A modern terminal task manager with extensive customization"
  homepage "https://github.com/spmfte/tuiodo"
  head "https://github.com/spmfte/tuiodo.git", branch: "master"
  version "1.1.1"
  license "MIT"

  # Additional metadata for the formula
  livecheck do
    url :stable
    strategy :github_latest
  end

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w")
  end

  def caveats
    <<~EOS
      TUIODO has been installed!
      
      Run `tuiodo` to start managing your tasks.
      For help and options, use `tuiodo --help`
      
      New in v1.1.1:
      - Enhanced task expansion UI with better formatting
      - Improved spacebar functionality for task completion
      - Circular cursor navigation (wraps around list edges)
      - Delete confirmation with undo capability
      - Advanced metadata tag support (@due, @tag, @status)
      - Color-coded progress bar in status line
    EOS
  end

  test do
    system "#{bin}/tuiodo", "--version"
  end
end 