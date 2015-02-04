Summary:    Export fields from Elasticsearch as tab-separated values.
Name:       estab
Version:    0.2.4
Release:    0
License:    MIT
BuildArch:  x86_64
BuildRoot:  %{_tmppath}/%{name}-build
Group:      System/Base
Vendor:     UB Leipzig
URL:        https://github.com/miku/estab

%description

Command line elasticsearch export tool. Exports TSV or raw documents.

%prep
# the set up macro unpacks the source bundle and changes in to the represented by
# %{name} which in this case would be my_maintenance_scripts. So your source bundle
# needs to have a top level directory inside called my_maintenance _scripts
# %setup -n %{name}

%build
# this section is empty for this example as we're not actually building anything

%install
# create directories where the files will be located
mkdir -p $RPM_BUILD_ROOT/usr/local/sbin

# put the files in to the relevant directories.
# the argument on -m is the permissions expressed as octal. (See chmod man page for details.)
install -m 755 estab $RPM_BUILD_ROOT/usr/local/sbin

%post
# the post section is where you can run commands after the rpm is installed.
# insserv /etc/init.d/my_maintenance

%clean
rm -rf $RPM_BUILD_ROOT
rm -rf %{_tmppath}/%{name}
rm -rf %{_topdir}/BUILD/%{name}

# list files owned by the package here
%files
%defattr(-,root,root)
/usr/local/sbin/estab


%changelog
* Wed Feb 4 2015 Martin Czygan
- 0.2.4 release
- add -precision flag

* Sat Jan 31 2015 Martin Czygan
- 0.2.3 release
- do not err on EOF

* Tue Jan 20 2015 Martin Czygan
- 0.2.2 release
- added -zero-as-null flag

* Thu Dec 18 2014 Martin Czygan
- 0.2.0 release
- added -raw flag
- use buffered output

* Wed Sep 16 2014 Martin Czygan
- 0.1.3 release
- add custom query support

* Wed Aug 21 2014 Martin Czygan
- 0.1.0 release
- initial release
