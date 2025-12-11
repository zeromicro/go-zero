{ pkgs ? import <nixpkgs> {} }:

pkgs.buildGoModule rec {
  name = "goctl";
  version = "1.8.3";

  src = pkgs.fetchFromGitHub {
    owner = "zeromicro";
    repo = "go-zero";
    tag = "v${version}";
    hash = "sha256-v5WzqMotF9C7i9hTYSjaPmTwveBVDVn+SKQXYuS4Rdc=";
  };

  vendorHash = "sha256-tOIlfYiAI9m7oTZyPDCzTXg9XTwBb6EOVLzDfZnzL4E=";

  modRoot = "tools/goctl";
  subPackages = [ "." ];

  ldflags = [
    "-s"
    "-w"
  ];

  meta = {
    description = "CLI handcuffle of go-zero, a cloud-native Go microservices framework";
    longDescription = ''
      goctl is a go-zero's built-in handcuffle that is a major
      lever to increase development efficiency, generating code,
      document, deploying k8s yaml, dockerfile, etc.
    '';
    license = pkgs.lib.licenses.mit;
    homepage = "https://go-zero.dev";
    mainProgram = "goctl";
  };
}
