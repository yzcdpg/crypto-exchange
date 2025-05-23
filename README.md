# crypto-exchange
数字货币交易所实践，主要用于学习交流。

## 市价单
市价单（Market Order）是一种以当前市场价格立即执行的交易指令，适用于快速买入或卖出资产。

### 基本定义
* 市价单是以当前市场最优价格立即成交的订单，不指定具体价格。
* 买入市价单：以当前市场最低卖出价（Ask Price，卖一价）或接近的价格买入。
* 卖出市价单：以当前市场最高买入价（Bid Price，买一价）或接近的价格卖出。

### 工作原理
#### 订单提交
* 交易者向交易所提交市价单，表明希望立即以市场价格买入或卖出一定数量的资产。
#### 价格撮合
* 交易所根据当前订单簿（Order Book）中的价格撮合订单。
* 买入市价单匹配订单簿中的最低卖出价（Ask）。
* 卖出市价单匹配订单簿中的最高买入价（Bid）。
* 如果订单量较大，可能会跨多个价格层级成交（例如，消耗多个卖单或买单）。 
#### 立即执行
* 市价单优先于限价单，交易所会尽快完成撮合，通常在几毫秒内执行。 
* 成交价格取决于订单簿的实时价格，可能略有波动。 
#### 成交确认
* 交易完成后，交易所返回成交价格和数量，订单状态更新为“已成交”。

### 示例
#### 买入市价单
* 当前订单簿：卖一价100元（1000股），卖二价101元（500股）。
* 你提交买入1200股的市价单：
* 1000股以100元成交。
* 剩余200股以101元成交。
* 平均成交价约为100.17元（假设无其他费用）。
#### 卖出市价单
* 当前订单簿：买一价99元（800股），买二价98元（600股）。
* 你提交卖出1000股的市价单：
* 800股以99元成交。
* 剩余200股以98元成交。
* 平均成交价约为98.8元。

### 注意事项
* 监控市场深度：在流动性低的资产或市场中，市价单可能导致严重滑点，建议查看订单簿。
* 避免极端波动：在市场剧烈波动时（如重大新闻发布），市价单可能以意外价格成交。
* 费用考虑：市价单通常作为Taker，费用较高，需了解交易所的费率结构。
* 替代选择：如果对价格敏感，可使用限价单（Limit Order）以控制成交价格，但可能无法立即成交。

## 限价单
限价单（Limit Order）是一种以指定价格或更优价格执行的交易指令，适用于希望控制成交价格的交易者

### 基本定义
* 限价单是交易者指定一个价格（限价），表示愿意以该价格或更优的价格买入或卖出资产。
  * 买入限价单：以指定价格或更低价格买入。
  * 卖出限价单：以指定价格或更高价格卖出。
* 限价单不会立即成交，除非市场价格达到或优于指定价格。


### 工作原理
#### 订单提交
* 交易者向交易所提交限价单，指定资产、数量和限价。
* 例如：买入100股，限价100元；或卖出100股，限价105元。
#### 加入订单簿
* 交易所将限价单放入订单簿（Order Book），按照价格优先、时间优先的原则排列。
* 买入限价单进入买单队列（Bid），卖出限价单进入卖单队列（Ask）。
#### 价格撮合
* 当市场价格达到或优于限价时，交易所撮合订单：
  * 买入限价单：当市场最低卖出价（Ask Price）低于或等于限价时成交。
  * 卖出限价单：当市场最高买入价（Bid Price）高于或等于限价时成交。
* 成交价格通常是限价或更优价格（例如，买入限价100元，可能以99元成交）。
#### 等待或取消
* 如果市场价格未达到限价，订单保持在订单簿中，直到：
  * 市场价格满足条件，订单成交。
  * 交易者手动取消订单。
  * 订单到期（取决于交易所设置的订单有效期，如当日有效或长期有效）。
#### 成交确认
* 成交后，交易所返回成交价格、数量和时间，订单状态更新为“已成交”或“部分成交”。

### 示例
#### 买入限价单：
* 当前市场：卖一价101元，买一价100元。
* 你提交买入100股，限价100元：
  * 订单进入买单队列，等待卖一价跌至100元或以下。
  * 如果卖一价降至100元，订单以100元成交；如果降至99元，可能以99元成交。
  * 如果价格始终高于100元，订单不会成交。
#### 卖出限价单：
* 当前市场：买一价99元，卖一价100元。
* 你提交卖出100股，限价100元：
  * 订单进入卖单队列，等待买一价涨至100元或以上。
  * 如果买一价升至100元，订单以100元成交；如果升至101元，可能以101元成交。
  * 如果价格始终低于100元，订单不会成交。

### 注意事项
* 市场波动：在高波动市场中，限价可能难以触及，导致错过交易机会。
* 订单簿深度：检查市场深度，确保限价合理，过高/低的限价可能长期未成交。
* 部分成交：大额订单可能分次成交，需关注订单状态。
* 时间优先：同一价格的限价单按提交时间排序，先提交者优先成交。
* 与市价单对比：限价单适合价格优先，市价单适合速度优先。

### 做多限价单
#### 原理
* 买入限价单：你指定一个价格（限价），表示你愿意以这个价格或更低的价格买入资产。
* 市价：当前市场上的最低卖出价（即卖一价，Ask Price）。

当你的买入限价高于或等于当前市场的最低卖出价时，交易所会将你的订单与卖单撮合，订单会立即执行。
如果你的限价高于当前市价（卖一价），意味着你愿意支付的价格高于市场上最便宜的卖单价格。
交易所会优先以当前市场的最低卖出价成交，而不是以你的限价成交（因为限价单保证价格不会高于你的设定价格）。
#### 例
当前市价（卖一价）是100元。
你设置的买入限价是105元。
你的订单会立即以100元（或接近市价的价格）成交，而不是105元，除非市场卖单迅速变动。

### 做空限价单
#### 原理
* 卖出限价单：你指定一个价格（限价），表示你愿意以这个价格或更高的价格卖出资产。
* 市价：当前市场上的最高买入价（即买一价，Bid Price）。

当你的卖出限价低于或等于当前市场的最高买入价时，交易所会将你的订单与买单撮合，订单会立即执行。
如果你的限价低于当前市价（买一价），意味着你愿意接受的价格低于市场上最高的买单价格。
交易所会优先以当前市场的最高买入价成交，而不是以你的限价成交（因为限价单保证价格不会低于你的设定价格）。

#### 例
当前市价（买一价）是100元。
你设置的卖出限价是95元。
你的订单会立即以100元（或接近市价的价格）成交，而不是95元，除非市场买单迅速变动。





