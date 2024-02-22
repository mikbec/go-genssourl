# go-genssourl

Der Dienst __go-genssourl__ ist ein Web-Dienst, der in der aktuellen Version
als Backend eines Web-Servers (wie bspw. Apache HTTPD oder NGINX) eingesetzt
werden kann, um anhand des durch den Web-Server authentifizierten Nutzers eine
Weiterleitungs-URL zu generieren, die folgenden Aufbau hat:

        <sever_protocol>://<server_host>[:<server_port>]/[<server_path>]?<url_attr_username_key>=username&<url_attr_timestamp_key>=timestamp&<url_attr_hash_key>=hash_val[&<url_attr_id_key>=<id_val>]

__Ein Beispiel:__ Wenn sich ein Nutzer

* als __user1__
* mit E-Mail-Adresse __user1@my.domain__
* um __2023-11-23 08:15:32 +00:00__

authentifiziert hat, dann koennten generierte URLs folgendermassen aussehen:

        http://server-one.my.domain/?user=user1&ts=2023-11-23T08%3A15%3A32Z&hash=12e3e5....34bc

oder

        https://server-two.my.domain:32146/issue?email=user1%40my.domain&timestamp=2023-11-23T08%3A15%3A32Z&key=12e3e5....34bc&id=id001

## URL-Parameter
Die Bedeutung der einzelnen URL-Parameter:

* __server_protocol__
    * das für die Weiterleitungs-URL zu verwendende Protokoll
    * Sinnvolle Werte sind hier _http_ oder _https_.
    * optionaler Parameter mit Standard-Einstellung _https_
* __server_host__
    * Name oder IP-Adresse des Zielservers
    * optionaler Parameter mit Standard-Einstellung _localhost_
* __server_port__
    * Port des Dienstes auf dem Zielserver
    * optionaler Parameter
* __server_path__
    * Pfad in der URL des Zielservers
    * optionaler Parameter
* __url_attr_username_key__
    * Name des Parameters, der den (durch den Web-Server) authentifizierten Accountnamen enthalten soll
    * optionaler Parameter mit Standard-Einstellung _user_
    * Format des Wertes
        * Es wird der Wert eingesetzt, der vom Web-Server im Authentifierungsvorgang ermittelt wurde.
        * Ggf. erfolgt ein URL-Encoding von Zeichen, die nicht in einer URL verwendet werden duerfen. Bspw. wird bei einer E-Mail-Adresse das __`@`__ Zeichen durch die HTML-Repraenstation __`%40`__ ersetzt.
* __url_attr_timestamp_key__
    * Name des Parameters, der den aktuellen Zeitstempel enthalten soll
    * optionaler Parameter mit Standard-Einstellung _ts_
    * Format des Wertes
        * Der Zeitstempel hat das Format als Layout-Wert __`2006-01-02T15:04:05Z`__
        * Ggf. erfolgt ein URL-Encoding von Zeichen, die nicht in einer URL verwendet werden duerfen.
        * Hinweis: Das Go Layout-Format besteht aus Teilen des Datums __`2006-01-02 15:04:05 -0700`__, der Zeitzone __`MST`__, der Tagesbezeichnungen __`Mon`__ oder __`Monday`__ u.s.w..
* __url_attr_hash_key__
    * Name des Parameters, der einen Hash-Wert enthalten soll, der aus dem Accountnamen und dem Zeitstempel berechnet wird.
    * optionaler Parameter mit Standard-Einstellung _hash_
    * Format des Wertes
        * Hexdezimalzahl des per RSA verschluesselten SHA1-Wertes aus Accountname und Zeitstempel
        * der Hash-Wert wird immer berechnet
* __url_attr_id_key__
    * Name des Parameters, der eine ID Schluessel enthalten soll
    * optionaler Parameter mit Standard-Einstellung _id_
    * Format des Wertes
        * Zeichenkette
        * falls diese Zeichenkette leer ist, wird das Attribut nicht gesetzt

## Die Berechnung des Hash-Wertes
Die Berechnung des Hash-Wertes erfolgt folgendermassen:

* Es wird eine Zeichenkette erzeugt, die aus dem Accountname und dem Zeitstempel besteht, welche ohne Trennzeichen zusammengefuegt werden.
* Von dieser Zeichenkette wird der SHA1-Hash-Wert berechnet.
* Dieser SHA1-Hash-Wert wird mit einem privaten RSA Schluessel verschluesselt.
* Dieser verschluesselte Byte-Wert wird in eine Hexadezimal-Repraesentation ueberfuehrt. 
