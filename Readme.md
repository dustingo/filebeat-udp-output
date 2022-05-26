### filebeat 的udp 输出插件
### 版本  
beats/V7
### 使用方法  
1. 手动注册到filebeat的main文件
```golang
package main

import (
	"os"

	"github.com/elastic/beats/v7/filebeat/cmd"
	inputs "github.com/elastic/beats/v7/filebeat/input/default-inputs"
    _ "gitee.com/kybeijing/filebeat-udp-output "
)

// The basic model of execution:
// - input: finds files in paths/globs to harvest, starts harvesters
// - harvester: reads a file, sends events to the spooler
// - spooler: buffers events until ready to flush to the publisher
// - publisher: writes to the network, notifies registrar
// - registrar: records positions of files read
// Finally, input uses the registrar information, on restart, to
// determine where in each file to restart a harvester.
func main() {
	if err := cmd.Filebeat(inputs.Init, cmd.FilebeatSettings()).Execute(); err != nil {
		os.Exit(1)
	}
}
```
2. 将代码copy到libeat的outputs，然后在main文件手动注册
```golang
import _ "github.com/elastic/beats/v7/libbeat/outputs/udp"
```
### 插件配置  
```yaml
output.udp:
  host: localhost
  port: 514
  bulk_max_size: 1000
  bulk_send_delay: 20
```
###用途