# Installation

## On Debian, Ubuntu, etc.

For Debian like systems do:

```bash
    # your are in the docs/ path of this project
    cp ../deploy/bin/genssourl /usr/sbin/go-genssourl
    cp go-genssourl.service.debian /lib/systemd/system/go-genssourl.service
    cp systemd-service.env /etc/default/go-genssourl
    
    mkdir /etc/go-genssourl
    cp config.yaml /etc/go-genssourl/config.yaml
    
    # run this command if you want to change something in the embedded content
    /usr/sbin/go-genssourl -webappdir /etc/go-genssourl/webapp -copyefs
    
    # now configure and start the service
    systemctl daemon-reload
    systemctl enable go-genssourl.service
    systemctl start go-genssourl.service 
    systemctl status go-genssourl.service
```

__Caution__: You will need a certificate for the hash calculation from the destination server admin.

## On OpenELA, CentOS, RHEL, OracleLinux, RockyLinux, AlmaLinux, etc.

For OpenELA like systems do:

```bash
    # your are in the docs/ path of this project
    cp ../deploy/bin/genssourl /usr/sbin/go-genssourl
    cp go-genssourl.service.openela /lib/systemd/system/go-genssourl.service
    cp systemd-service.env /etc/sysconfig/go-genssourl
    
    mkdir /etc/go-genssourl
    cp config.yaml /etc/go-genssourl/config.yaml
    
    # run this command if you want to change something in the embedded content
    /usr/sbin/go-genssourl -webappdir /etc/go-genssourl/webapp -copyefs
    
    # now configure and start the service
    systemctl daemon-reload
    systemctl enable go-genssourl.service
    systemctl start go-genssourl.service 
    systemctl status go-genssourl.service
```

__Caution__: You will need a certificate for the hash calculation from the destination server admin.
