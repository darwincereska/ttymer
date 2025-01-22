{ lib
, buildGoModule
, fetchFromGitHub
}:

buildGoModule {
  pname = "ttymer";
  version = "1.0.1";

  src = fetchFromGitHub {
    owner = "darwincereska";
    repo = "ttymer";
    rev = "v1.0.1";
    hash = "sha256-a3+TAGBz1br2TCu9FxtUN4G3H84NZwwV/fFg5/HbJ2k=";
  };

  vendorHash = null;

  meta = with lib; {
    description = "Terminal based timer";
    homepage = "https://github.com/darwincereska/ttymer";
    license = licenses.mit;
    maintainers = [ maintainers.darwincereska ];
  };
}