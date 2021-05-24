CREATE TABLE tbl_acct_config(
  acctid INT PRIMARY KEY auto_increment  NOT NULL, 
  cellphone VARCHAR(64) NOT NULL UNIQUE, 
  realname VARCHAR(64) NOT NULL, 
  idcard VARCHAR(64) NOT NULL UNIQUE, 
  pubkey VARCHAR(512),
  role INT NOT NULL, 
  state INT NOT NULL, 
  createtime DATETIME, 
  updatetime DATETIME);

CREATE TABLE tbl_acct_wallet_relation(
  relationid INT PRIMARY KEY auto_increment  NOT NULL, 
  acctid INT NOT NULL, 
  walletid INT NOT NULL, 
  createtime DATETIME);

CREATE TABLE tbl_coin_config(
  coinid INT PRIMARY KEY auto_increment  NOT NULL, 
  coinsymbol VARCHAR(16) NOT NULL UNIQUE, 
  ip VARCHAR(64) NOT NULL, 
  rpcport INT NOT NULL, 
  rpcuser VARCHAR(64), 
  rpcpass VARCHAR(64), 
  state INT NOT NULL, 
  createtime DATETIME, 
  updatetime DATETIME);

CREATE TABLE tbl_notification(
  notifyid INT PRIMARY KEY auto_increment  NOT NULL, 
  acctid INT NOT NULL, 
  wallettid INT, 
  trxid INT, 
  notifytype INT NOT NULL, 
  notification TEXT, 
  state INT NOT NULL, 
  reserved1 TEXT, 
  reserved2 TEXT, 
  createtime DATETIME, 
  updatetime DATETIME);

CREATE TABLE tbl_operator_log(
  logid INT PRIMARY KEY auto_increment  NOT NULL, 
  acctid INT NOT NULL, 
  optype INT NOT NULL, 
  content TEXT, 
  createtime DATETIME);

CREATE TABLE tbl_pubkey_pool(
  keyindex INT PRIMARY KEY NOT NULL, 
  pubkey VARCHAR(512) NOT NULL, 
  isused BOOL NOT NULL, 
  createtime DATETIME, 
  usedtime DATETIME);

CREATE TABLE tbl_sequence(
  seqvalue INT PRIMARY KEY auto_increment  NOT NULL, 
  idtype INT NOT NULL, 
  state INT NOT NULL, 
  createtime DATETIME NOT NULL, 
  updatetime DATETIME);

CREATE TABLE tbl_transaction(
  trxid INT PRIMARY KEY auto_increment  NOT NULL, 
  rawtrxid VARCHAR(128), 
  walletid INT NOT NULL, 
  coinid INT NOT NULL, 
  contractaddr VARCHAR(128), 
  acctid INT NOT NULL, 
  fromaddr VARCHAR(128) NOT NULL, 
  toaddr VARCHAR(128) NOT NULL, 
  amount VARCHAR(128) NOT NULL, 
  feecost VARCHAR(128), 
  trxtime DATETIME, 
  needconfirm INT NOT NULL, 
  confirmed INT NOT NULL, 
  acctconfirmed TEXT NOT NULL, 
  fee VARCHAR(128), 
  gasprice VARCHAR(128), 
  gaslimit VARCHAR(128), 
  state INT NOT NULL,
  signature VARCHAR(512) NOT NULL);

CREATE TABLE tbl_wallet_config(
  walletid INT PRIMARY KEY auto_increment  NOT NULL, 
  coinid INT(16) NOT NULL, 
  walletname VARCHAR(64) NOT NULL UNIQUE, 
  keyindex INT NOT NULL UNIQUE, 
  address VARCHAR(64) NOT NULL UNIQUE, 
  destaddress TEXT, 
  needsigcount INT NOT NULL, 
  fee VARCHAR(64), 
  gasprice VARCHAR(64), 
  gaslimit VARCHAR(64), 
  state INT NOT NULL, 
  createtime DATETIME, 
  updatetime DATETIME);

