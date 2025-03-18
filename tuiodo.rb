class Tuiodo < Formula
  desc "A modern terminal task manager with extensive customization"
  homepage "https://github.com/spmfte/tuiodo"
  url "https://github.com/spmfte/tuiodo/archive/refs/tags/v1.1.2.tar.gz"
  sha256 "6186659032fa47685499ad52b59f9e4ea1c9c29dc18197b761d9497d46fe3581"
  license "MIT"
  head "https://github.com/spmfte/tuiodo.git", branch: "master"

  # Additional metadata for the formula
  livecheck do
    url :stable
    strategy :github_latest
  end

  depends_on "go" => :build

  def install
    # Set build time, version and commit hash
    if build.head?
      # For HEAD installs, get the actual git commit
      git_revision = Utils.git_short_head
      git_version = "head"
    else
      # For release installs, use the version
      git_revision = stable.version.to_s[0..6]
      git_version = version
    end

    ldflags = %W[
      -s -w
      -X main.Version=#{git_version}
      -X main.BuildTime=#{Time.now.iso8601}
      -X main.GitCommit=#{git_revision}
    ].join(" ")

    system "go", "build", *std_go_args(ldflags: ldflags)
  end

  def caveats
    <<~EOS
      TUIODO has been installed!
      
      Run `tuiodo` to start managing your tasks.
      For help and options, use `tuiodo --help`
      
      New in v1.1.2:
      - Customizable category colors via configuration
      - Default sort by priority for better task organization
      - Categories for completed tasks now appear with dimmed colors
      - Completed tasks always appear at the bottom of lists
      - Priority indicators hidden for completed tasks
    EOS
  end

  test do
    system "#{bin}/tuiodo", "--version"
  end
end