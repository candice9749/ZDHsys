# -*- coding: cp936 -*-
import xmlrpclib


def cobbler(name, hostname, profile, ks_meta, modify_interface):
    try:
        server = xmlrpclib.Server("http://localhost/cobbler_api")
        token = server.login("cobbler", "cobbler")
        system_id = server.new_system(token)

        server.modify_system(system_id, "name", name, token)
        server.modify_system(system_id, "kernel_options", "net.ifnames=0 biosdevname=0", token)
        server.modify_system(system_id, "hostname", hostname, token)
        server.modify_system(system_id, "profile", profile, token)
        server.modify_system(system_id, "ks_meta", ks_meta, token)
        try:
            print modify_interface
            interface = eval(modify_interface)
        except Exception as e:
            print "kkkkkkk",e
            return
        server.modify_system(system_id, 'modify_interface', interface, token)
        server.save_system(system_id, token)
        server.sync(token)
    except Exception as e:
        print "kkkkkkk",e