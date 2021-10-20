%{!?_version: %define _version 0.1 }
%{!?_major: %define _major %{_version} }
%{!?_branch: %define _branch %(git rev-parse --abbrev-ref HEAD) }

%define _release %(/bin/date +"%Y%m%d.%H%M")

Name:           scli
Version:        %{_version}
Release:        %{_release}
Source0:        %{name}-%{_major}.tar.gz
BuildRoot:      %{_tmppath}/%{name}
License:        Proprietary
Vendor:         maxharitonov
Group:          cli
Summary:        sample cli

%description
sample cli

Packager: %(echo ${USER})
BuildHost: %(echo ${HOSTNAME})
Branch: %{_branch}
Commit: %(git rev-parse HEAD)

%prep
[ "%{buildroot}" != "/" ] && rm -rf %{buildroot}
%setup -q -n %{name}
mkdir configs/tmp
go run cmd/scli/main.go completion bash > configs/tmp/bash

%build
go build -o build -ldflags "-linkmode=external \
    	-X gitlab.corp.mail.ru/cloud/go-libs/version.BuildDate=`date -u --rfc-3339=seconds | sed -e 's/ /T/'` \
    	-X gitlab.corp.mail.ru/cloud/go-libs/version.GitCommit=%(git rev-parse HEAD) \
    	-X gitlab.corp.mail.ru/cloud/go-libs/version.GitBranch=%{_branch} \
    	-X gitlab.corp.mail.ru/cloud/go-libs/version.BuildClean=`if [ -n "%(git status --porcelain)" ]; then echo 'no'; else echo 'yes'; fi` \
    	-X gitlab.corp.mail.ru/cloud/go-libs/version.BuildVersion=%(git describe --tags --dirty --always|sed 's/-/./')" \
		./cmd/...

%install
[ "%{buildroot}" != "/" ] && rm -rf %{buildroot}
%{__mkdir} -p %{buildroot}/usr/bin
%{__mkdir} -p %{buildroot}/etc/%{name}

%{__install} -pD -m 755 build/scli %{buildroot}/usr/bin
%{__install} -pD -m 644 configs/config.yaml  %{buildroot}/etc/%{name}/config.yaml

%{__install} -pD -m 644 configs/tmp/bash %{buildroot}/etc/bash_completion.d/%{name}

%clean
[ "%{buildroot}" != "/" ] && rm -rf %{buildroot}
rm -rf configs/tmp

%files
/usr/bin/%{name}
%config(noreplace) /etc/%{name}/config.yaml
/etc/bash_completion.d/%{name}
