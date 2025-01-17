package provider

import (
	"errors"
	"github.com/GarnerCorp/go-zabbix-api"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceZabbixProxy() *schema.Resource {
	return &schema.Resource{
		Create: resourceZabbixProxyCreate,
		Read:   resourceZabbixProxyRead,
		Update: resourceZabbixProxyUpdate,
		Delete: resourceZabbixProxyDelete,
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Technical name of the proxy.",
			},
			"proxyid": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "(readonly) ID of the proxy.",
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"interfaces": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     interfaceSchema,
				Optional: true,
			},
			"tls_connect": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"tls_accept": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"tls_psk_identity": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tls_psk": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"generate_tls_psk": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"passive": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
		},
	}
}

func createProxyObj(d *schema.ResourceData, api *zabbix.API) (*zabbix.Proxy, error) {

	proxy := zabbix.Proxy{
		Host:           d.Get("host").(string),
		Description:    d.Get("description").(string),
		TlsConnect:     d.Get("tls_connect").(int),
		TlsAccept:      d.Get("tls_accept").(int),
		TlsPskIdentity: d.Get("tls_psk_identity").(string),
		TlsPsk:         d.Get("tls_psk").(string),

		Status: 5,
	}

	if d.Get("generate_tls_psk").(bool) {
		if proxy.TlsPsk == "" {
			randomPsk, err := randomHex(32)

			if err != nil {
				return nil, err
			}

			proxy.TlsPsk = randomPsk
		}
	}

	//5 is active, 6 - passive proxy
	if d.Get("passive").(bool) {
		proxy.Status = 6

		interfaces, err := getInterfaces(d)

		if err != nil {
			return nil, err
		}

		if len(interfaces) <= 0 {
			return &proxy, errors.New("passive proxies require an interface")
		}

		proxy.Interfaces = interfaces
		proxy.Interface = interfaces[0]
	}

	return &proxy, nil
}

func resourceZabbixProxyCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	proxy, err := createProxyObj(d, api)

	if err != nil {
		return err
	}

	proxys := zabbix.Proxies{*proxy}

	err = api.ProxiesCreate(proxys)

	if err != nil {
		return err
	}

	log.Printf("Created proxy id is %s", proxys[0].ProxyId)

	d.Set("proxyid", proxys[0].ProxyId)
	d.SetId(proxys[0].ProxyId)

	return nil
}

func resourceZabbixProxyRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	log.Printf("Will read proxy with id %s", d.Get("proxyid").(string))

	proxy, err := api.ProxyGetById(d.Get("proxyid").(string))

	if err != nil {
		return err
	}

	d.Set("host", proxy.Host)
	d.Set("description", proxy.Description)

	d.Set("tls_connect", proxy.TlsConnect)
	d.Set("tls_accept", proxy.TlsAccept)
	d.Set("tls_psk", proxy.TlsPsk)
	d.Set("tls_psk_identity", proxy.TlsPskIdentity)

	d.Set("passive", proxy.Status == 5)

	return nil
}

func resourceZabbixProxyUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	proxy, err := createProxyObj(d, api)

	if err != nil {
		return err
	}

	proxy.ProxyId = d.Id()

	//interfaces can't be updated, changes will trigger recreate
	//sending previous values will also fail the update
	proxy.Interfaces = nil

	proxys := zabbix.Proxies{*proxy}

	err = api.ProxiesUpdate(proxys)

	if err != nil {
		return err
	}

	log.Printf("Created proxy id is %s", proxys[0].ProxyId)

	return nil
}

func resourceZabbixProxyDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	return api.ProxiesDeleteByIds([]string{d.Id()})
}
