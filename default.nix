{ lib, buildGoModule }:
buildGoModule {
  pname = "pad-archiver";
  version = "0.0.1";
  src = ../.;
  vendorHash = "sha256-Z0Kxw0hX4z5NuXJNqtuRw6ZItCaXnTX44Vpv9IFaS38=";
  ldflags = [ "-s" "-w" ];
  meta = with lib; {
    description = "Archives Etherpads with git";
    homepage = "https://github.com/dadada/pad-archiver";
    license = licenses.mit;
    maintainers = [ "dadada" ];
  };
}
