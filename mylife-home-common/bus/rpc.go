package bus

import (
	"fmt"
	"mylife-home-common/tools"
)

const rpcDomain = "rpc"
const rpcServices = "services"
const rpcReplies = "replies"

const RpcTimeout = 2000

/*
class RpcError extends Error {
  constructor(public readonly remoteMessage: string, public readonly remoteStacktrace: string) {
    super(`A remote error occured: ${remoteMessage}`);
  }
}
*/

type Rpc struct {
	client   *client
	services map[string]RpcService
}

func newRpc(client *client) *Rpc {
	return &Rpc{
		client:   client,
		services: make(map[string]RpcService),
	}
}

type RpcService interface {
	init(client *client, address string) error
	terminate() error
}

func (rpc *Rpc) Serve(address string, svc RpcService) error {
	_, exists := rpc.services[address]
	if exists {
		panic(fmt.Errorf("service with address '%s' does already exist", address))
	}

	if err := svc.init(rpc.client, address); err != nil {
		return err
	}

	rpc.services[address] = svc
	return nil
}

func (rpc *Rpc) Unserve(address string) error {
	svc, exists := rpc.services[address]
	if !exists {
		panic(fmt.Errorf("service with address '%s' does not exist", address))
	}

	err := svc.terminate()
	delete(rpc.services, address)
	return err
}

/*
func (rpc *Rpc) Call(targetInstance string, address string, data any, timeout int) (any, error) {
	const id = randomTopicPart();
	const replyTopic = this.client.buildTopic(DOMAIN, REPLIES, id);
	const request: Request = { input: data, replyTopic };
	let buffer: Buffer;

	const messageWaiter = new MessageWaiter(this.client, address, replyTopic);
	await messageWaiter.init();
	try {
		await this.client.publish(this.client.buildRemoteTopic(targetInstance, DOMAIN, SERVICES, address), encoding.writeJson(request));
		buffer = await messageWaiter.waitForMessage(timeout);
	} finally {
		await messageWaiter.terminate();
	}

	const response = encoding.readJson(buffer) as Response;
	const { error, output } = response;
	if (error) {
		throw new RpcError(error.message, error.stacktrace);
	}

	return output;
}
*/

type rpcServiceImpl[TInput any, TOutput any] struct {
	client         *client
	address        string
	implementation func(TInput) (TOutput, error)
	msgToken       tools.RegistrationToken
}

func NewRpcService[TInput any, TOutput any](implementation func(TInput) (TOutput, error)) RpcService {
	return &rpcServiceImpl[TInput, TOutput]{
		implementation: implementation,
	}
}

func (svc *rpcServiceImpl[TInput, TOutput]) init(client *client, address string) error {
	svc.client = client
	svc.address = address

	svc.msgToken = svc.client.OnMessage().Register(svc.onMessage)
	return svc.client.Subscribe(svc.buildTopic())
}

func (svc *rpcServiceImpl[TInput, TOutput]) terminate() error {
	err := svc.client.Unsubscribe(svc.buildTopic())
	svc.client.OnMessage().Unregister(svc.msgToken)
	return err
}

func (svc *rpcServiceImpl[TInput, TOutput]) onMessage(m *message) {
	if m.InstanceName() != svc.client.InstanceName() || m.Domain() != rpcServices || m.Path() != svc.address {
		return
	}

	fireAndForget(func() error {
		var req request[TInput]
		Encoding.ReadTypedJson(m.Payload(), &req)
		resp := svc.handle(&req)
		output := Encoding.WriteJson(resp)
		return svc.client.Publish(req.ReplyTopic, output, false)
	})
}

func (svc *rpcServiceImpl[TInput, TOutput]) handle(req *request[TInput]) *response[TOutput] {
	output, err := svc.implementation(req.Input)

	if err != nil {
		return &response[TOutput]{
			Error: &reponseError{
				Message:    err.Error(),
				Stacktrace: tools.GetStackTraceStr(err),
			},
		}
	} else {
		return &response[TOutput]{
			Output: &output,
		}
	}
}

func (svc *rpcServiceImpl[TInput, TOutput]) buildTopic() string {
	return svc.client.BuildTopic(rpcDomain, rpcServices, svc.address)
}

/*
class MessageWaiter {
  constructor(private readonly client: Client, private readonly callAddress: string, private readonly topic: string) {
  }

  async init() {
    await this.client.subscribe(this.topic);
  }

  async waitForMessage(timeout: number) {
    return await new Promise<Buffer>((resolve, reject) => {

      const onEnd = () => {
        clearTimeout(timer);
        this.client.off('message', messageCb);
      };

      const messageCb = (mtopic: string, payload: Buffer) => {
        if (this.topic !== mtopic) {
          return;
        }

        onEnd();
        resolve(payload);
      };

      const timer = setTimeout(() => {
        onEnd();
        reject(new Error(`Timeout occured while waiting for message on topic '${this.topic}' (call address: '${this.callAddress}', timeout: ${timeout})`));
      }, timeout);

      this.client.on('message', messageCb);
    });
  }

  async terminate() {
    await this.client.unsubscribe(this.topic);
  }
}

const CHARSET = 'ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
const CHARSET_LEN = CHARSET.length;
const LEN = 16;

function randomTopicPart() {
  const array = new Array(LEN);
  for (let i = 0; i < LEN; ++i) {
    array[i] = CHARSET.charAt(Math.floor(Math.random() * CHARSET_LEN));
  }
  return array.join('');
}
*/

type request[TInput any] struct {
	Input      TInput `json:"input"`
	ReplyTopic string `json:"replyTopic"`
}

type response[TOutput any] struct {
	Output *TOutput      `json:"output"`
	Error  *reponseError `json:"error"`
}

type reponseError struct {
	Message    string `json:"message"`
	Stacktrace string `json:"stacktrace"`
}
