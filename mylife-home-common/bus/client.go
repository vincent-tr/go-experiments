package bus

import (
	"mylife-home-common/config"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type busConfig struct {
	ServerUrl string `mapstructure:"serverUrl"`
}

type Client struct {
	instanceName string
	mqtt         mqtt.Client
}

func NewClient(instanceName string) *Client {
	conf := busConfig{}
	config.BindStructure("bus", &conf)

	// Need it for topics
	client := &Client{
		instanceName: instanceName,
	}

	options := mqtt.NewClientOptions()
	options.AddBroker(conf.ServerUrl)
	options.SetClientID(instanceName)
	options.SetCleanSession(true)
	options.SetBinaryWill(client.buildTopic("online"), []byte{}, 0, true)

	client.mqtt = mqtt.NewClient(options)

	// client.mqtt.Connect()
	/*
	   	this.client.on('connect', () =>
	   	fireAsync(async () => {
	   		// given the spec, it is unclear if LWT should be executed in case of client takeover, so we run it to be sure
	   		await this.clearRetain(this.buildTopic('online'));

	   		await this.clearResidentState();
	   		await this.publish(this.buildTopic('online'), encoding.writeBool(true), true);
	   		this.onlineChange(true);

	   		if (this.subscriptions.size) {
	   			await this.client.subscribe(Array.from(this.subscriptions));
	   		}
	   	})
	   );

	   this.client.on('close', () => this.onlineChange(false));

	   this.client.on('error', (err) => {
	   	log.error(err, 'mqtt error');
	   });

	   this.client.on('message', (topic, payload) => this.emit('message', topic, payload));
	*/
	//	mqttClient.Publish()
	//	mqttClient.Subscribe()

	return client
}

/*


  private async clearResidentState() {
    // register on self state, and remove on every message received
    // wait 1 sec after last message receive
    const { promise: sleepPromise, reset: resetSleep } = sleepWithReset(this.residentStateDelay);

    const clearTopic = (topic: string, payload: Buffer, packet: mqtt.IPublishPacket) => {
      // only clear real retained messages
      if (packet.retain && payload.length > 0 && topic.startsWith(this.instanceName + '/')) {
        resetSleep();
        fireAsync(() => this.clearRetain(topic));
      }
    };

    const selfStateTopic = this.buildTopic('#');
    this.client.on('message', clearTopic);
    await this.subscribe(selfStateTopic);

    await sleepPromise;

    this.client.off('message', clearTopic);
    await this.unsubscribe(selfStateTopic);
  }

  private onlineChange(value: boolean): void {
    if (value === this._online) {
      return;
    }
    this._online = value;
    log.info(`online: ${value}`);
    this.emit('onlineChange', value);
  }

  get online(): boolean {
    return this._online;
  }

  async terminate(): Promise<void> {
    if (this.client.connected) {
      await this.clearRetain(this.buildTopic('online'));
      await this.clearResidentState();
    }
    await this.client.end(true);
  }

  async publish(topic: string, payload: Buffer, retain: boolean = false) {
    await this.client.publish(topic, payload, { retain });
  }

  async subscribe(topic: string | string[]) {
    if (!Array.isArray(topic)) {
      topic = [topic];
    }
    for (const item of topic) {
      this.subscriptions.add(item);
    }
    if (this.online) {
      await this.client.subscribe(topic);
    }
  }

  async unsubscribe(topic: string | string[]) {
    if (!Array.isArray(topic)) {
      topic = [topic];
    }
    for (const item of topic) {
      this.subscriptions.delete(item);
    }
    if (this.online) {
      await this.client.unsubscribe(topic);
    }
  }
}

function sleepWithReset(delay: number) {
  const deferred = new Deferred<void>();
  let timeoutHandle: NodeJS.Timeout;

  const reset = () => {
    clearTimeout(timeoutHandle);
    timeoutHandle = setTimeout(() => deferred.resolve(), delay);
  };

  reset();

  return { promise: deferred.promise, reset };
}

*/

func (client *Client) buildTopic(domain string, args ...string) string {
	finalArgs := append([]string{client.instanceName, domain}, args...)
	return strings.Join(finalArgs, "/")
}

func (client *Client) buildRemoteTopic(targetInstance string, domain string, args ...string) string {
	finalArgs := append([]string{targetInstance, domain}, args...)
	return strings.Join(finalArgs, "/")
}

func (client *Client) clearRetain(topic string) mqtt.Token {
	return client.mqtt.Publish(topic, 0, true, []byte{})
}
