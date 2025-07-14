Name:           restic-exporter
Version:        0.1.5
Release:        1%{?dist}
Summary:        Exports Prometheus metrics for a restic repository

License:        GPLv2
URL:            https://github.com/Rherer/restic-exporter
Source0:        https://github.com/Rherer/restic-exporter/archive/v%version.tar.gz       

BuildRequires: golang
BuildRequires: systemd-rpm-macros  
Requires: restic

%description
Exports metrics for a restic repository.
Settings can be configured using environment Variables.
For more see: https://github.com/Rherer/restic-exporter

%prep
%autosetup

%build
go build -a -ldflags "-B 0x$(head -c20 /dev/urandom|od -An -tx1|tr -d ' \n')" -v -o %{name}

%install
install -Dpm 0755 %{name} %{buildroot}%{_bindir}/%{name}
install -Dpm 644 %{name}.service %{buildroot}%{_unitdir}/%{name}.service

%post
%systemd_post %{name}.service

%preun
%systemd_preun %{name}.service

%files
%{_bindir}/%{name}
%{_unitdir}/%{name}.service


%changelog
* Thu Jun 12 2025 Raphael Eherer <raphael.eherer@gmail.com>
- 
