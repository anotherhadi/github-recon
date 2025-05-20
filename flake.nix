{
  description =
    "GH-Recon: Fetches and aggregates public OSINT data for a GitHub user, leveraging Go and the GitHub API.";

  inputs = { nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable"; };

  outputs = { self, nixpkgs }:
    let
      supportedSystems = [ "x86_64-linux" "aarch64-linux" ];

      forAllSystems = f:
        nixpkgs.lib.genAttrs supportedSystems
        (system: f system (import nixpkgs { inherit system; }));

      pname = "gh-recon";
      version = "0.2.1";

      ldflags = [ "-s" "-w" ];

    in {
      packages = forAllSystems (system: pkgs: {
        "${pname}" = pkgs.buildGoModule {
          inherit pname version ldflags;

          src = ./.;

          vendorHash = "sha256-S8IzmdiVvBtnQQl0AewGZ1yuitvrdnVQ/Jf2230g3Mg=";

          meta = with pkgs.lib; {
            description =
              "Fetches and aggregates public OSINT data for a GitHub user.";
            homepage = "https://github.com/anotherhadi/gh-recon";
            platforms = platforms.unix;
          };
        };
      });

      defaultPackage =
        forAllSystems (system: pkgs: self.packages.${system}.${pname});
    };
}
