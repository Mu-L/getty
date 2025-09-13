#!/usr/bin/env bash
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# ******************************************************
# DESC    :
# AUTHOR  : Alex Stocks
# VERSION : 1.0
# LICENCE : Apache License 2.0
# EMAIL   : alexstocks@foxmail.com
# MOD     : 2019-09-07 18:32
# FILE    : build_all.sh
# ******************************************************


cd echo/tcp-echo/client && sh assembly/mac/dev.sh && rm -rf target && cd -
cd echo/tcp-echo/server && sh assembly/mac/dev.sh && rm -rf target && cd -

cd echo/udp-echo/client && sh assembly/mac/dev.sh && rm -rf target && cd -
cd echo/udp-echo/server && sh assembly/mac/dev.sh && rm -rf target && cd -

cd echo/ws-echo/client && sh assembly/mac/dev.sh && rm -rf target && cd -
cd echo/ws-echo/server && sh assembly/mac/dev.sh && rm -rf target && cd -

cd echo/wss-echo/client && sh assembly/mac/dev.sh && rm -rf target && cd -
cd echo/wss-echo/server && sh assembly/mac/dev.sh && rm -rf target && cd -

# cd rpc/client && sh assembly/mac/dev.sh && rm -rf target && cd -
# cd rpc/server && sh assembly/mac/dev.sh && rm -rf target && cd -

# cd micro/client && sh assembly/mac/dev.sh && rm -rf target && cd -
# cd micro/server && sh assembly/mac/dev.sh && rm -rf target && cd -

