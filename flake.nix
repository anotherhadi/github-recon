{
  description = "Retrieves and aggregates public OSINT data about a Github user using Go and the Github API. Finds hidden emails in commit history, previous usernames, friends, other Github accounts, and more.";

  inputs = {nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";};

  outputs = {
    self,
    nixpkgs,
  }: let
    supportedSystems = ["x86_64-linux" "aarch64-linux"];

    forAllSystems = f:
      nixpkgs.lib.genAttrs supportedSystems
      (system: f system (import nixpkgs {inherit system;}));

    pname = "github-recon";
    version = "2.0.0";

    ldflags = ["-s" "-w"];
  in {
    packages = forAllSystems (system: pkgs: {
      "${pname}" = pkgs.buildGoModule {
        inherit pname version ldflags;

        src = ./.;
        subPackages = ["cmd"];
        outputs = ["out"];
        installPhase = ''
          mkdir -p $out/bin
          cp $GOPATH/bin/cmd $out/bin/github-recon
        '';

        vendorHash = "sha256-AD0h0k2n8gPqSBz5qqb0ZON/jWiSEWpeO97xR7cYSy8=";

        meta = with pkgs.lib; {
          description = "Retrieves and aggregates public OSINT data about a Github user using Go and the Github API. Finds hidden emails in commit history, previous usernames, friends, other Github accounts, and more.";
          homepage = "https://github.com/anotherhadi/github-recon";
          platforms = platforms.unix;
        };
      };
    });

    defaultPackage =
      forAllSystems (system: pkgs: self.packages.${system}.${pname});
  };
}
