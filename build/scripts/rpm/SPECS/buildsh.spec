Name:           %{_product_name}
Version:        %{_product_version}

Release:        1.el%{_rhel_version}
Summary:        Buildsh is docker powered shell that makes it easy to run a script
Group:          Development/Tools
License:        MIT
Source0:        %{name}_linux_amd64.zip
BuildRoot:      %(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)

%description
Buildsh is docker powered shell that makes it easy to run a script

%prep
%setup -q -c

%install
mkdir -p %{buildroot}/%{_bindir}
cp %{name} %{buildroot}/%{_bindir}

%pre

%post

%preun

%clean
rm -rf %{buildroot}


%files
%defattr(-,root,root,-)
%attr(755, root, root) %{_bindir}/%{name}

%doc
