export const proto = `
openapi: 3.0.0
info:
  version: 0.0.3
  title: ai-proto
  description: ai proto for coder
  contact:
    name: linksaas
    email: panleiming@linksaas.pro
    url: https://github.com/linksaas/ai-proto
servers:
  - url: __SERVER_ADDR__
tags:
  - name: dev
    description: 开发模式相关接口
  - name: coding
    description: 提高编码效率的工具接口
paths:
  /api/coding/complete/{lang}:
    post:
      tags:
        - coding
      summary: 根据上下文补全代码
      description: 根据上下文补全代码
      parameters:
        - $ref: '#/components/parameters/AuthToken'
        - $ref: '#/components/parameters/Lang'
      requestBody:
        required: true
        content:
          text/plain:
            schema:
              type: string
      responses:
        '200':
          description: 成功
          headers:
            Access-Control-Allow-Origin:
              schema:
                type: string
                default: '*'
          content:
            application/json:
              schema:
                type: array
                items:
                  type: string
                description: 生成代码候选列表
        '403':
          description: 没有权限
          headers:
            Access-Control-Allow-Origin:
              schema:
                type: string
                default: '*'
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrInfo'
        '500':
          description: 失败
          headers:
            Access-Control-Allow-Origin:
              schema:
                type: string
                default: '*'
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrInfo'
  /api/coding/convert/{lang}:
    post:
      tags:
        - coding
      summary: 对选中代码转换成其他编程语言
      description: 对选中代码转换成其他编程语言
      parameters:
        - $ref: '#/components/parameters/AuthToken'
        - $ref: '#/components/parameters/Lang'
        - $ref: '#/components/parameters/DestLang'
      requestBody:
        required: true
        content:
          text/plain:
            schema:
              type: string
      responses:
        '200':
          description: 成功
          headers:
            Access-Control-Allow-Origin:
              schema:
                type: string
                default: '*'
          content:
            application/json:
              schema:
                type: array
                items:
                  type: string
                description: 生成代码候选列表
        '403':
          description: 没有权限
          headers:
            Access-Control-Allow-Origin:
              schema:
                type: string
                default: '*'
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrInfo'
        '500':
          description: 失败
          headers:
            Access-Control-Allow-Origin:
              schema:
                type: string
                default: '*'
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrInfo'
  /api/coding/explain/{lang}:
    post:
      tags:
        - coding
      summary: 解释选择代码
      description: 解释选择代码
      parameters:
        - $ref: '#/components/parameters/AuthToken'
        - $ref: '#/components/parameters/Lang'
      requestBody:
        required: true
        content:
          text/plain:
            schema:
              type: string
      responses:
        '200':
          description: 成功
          headers:
            Access-Control-Allow-Origin:
              schema:
                type: string
                default: '*'
          content:
            application/json:
              schema:
                type: array
                items:
                  type: string
                description: 生成解释结果候选列表
        '403':
          description: 没有权限
          headers:
            Access-Control-Allow-Origin:
              schema:
                type: string
                default: '*'
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrInfo'
        '500':
          description: 失败
          headers:
            Access-Control-Allow-Origin:
              schema:
                type: string
                default: '*'
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrInfo'
  /api/coding/fixError/{lang}:
    post:
      tags:
        - coding
      summary: 根据错误提示给出解决方案
      description: 根据错误提示给出解决方案
      parameters:
        - $ref: '#/components/parameters/AuthToken'
        - $ref: '#/components/parameters/Lang'
      requestBody:
        required: true
        content:
          text/plain:
            schema:
              type: string
      responses:
        '200':
          description: 成功
          headers:
            Access-Control-Allow-Origin:
              schema:
                type: string
                default: '*'
          content:
            application/json:
              schema:
                type: array
                items:
                  type: string
                description: 修复方案候选列表
        '403':
          description: 没有权限
          headers:
            Access-Control-Allow-Origin:
              schema:
                type: string
                default: '*'
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrInfo'
        '500':
          description: 失败
          headers:
            Access-Control-Allow-Origin:
              schema:
                type: string
                default: '*'
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrInfo'
  /api/coding/genTest/{lang}:
    post:
      tags:
        - coding
      summary: 对选中函数生成测试代码
      description: 对选中函数生成测试代码
      parameters:
        - $ref: '#/components/parameters/AuthToken'
        - $ref: '#/components/parameters/Lang'
      requestBody:
        required: true
        content:
          text/plain:
            schema:
              type: string
      responses:
        '200':
          description: 成功
          headers:
            Access-Control-Allow-Origin:
              schema:
                type: string
                default: '*'
          content:
            application/json:
              schema:
                type: array
                items:
                  type: string
                description: 生成代码候选列表
        '403':
          description: 没有权限
          headers:
            Access-Control-Allow-Origin:
              schema:
                type: string
                default: '*'
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrInfo'
        '500':
          description: 失败
          headers:
            Access-Control-Allow-Origin:
              schema:
                type: string
                default: '*'
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrInfo'
  /api/dev/cap:
    post:
      tags:
        - dev
      summary: 获取ai能力列表
      description: 获取ai能力列表
      parameters:
        - $ref: '#/components/parameters/AuthToken'
      requestBody:
        required: false
        content:
          application/x-www-form-urlencoded:
            schema:
              type: object
      responses:
        '200':
          description: 成功
          content:
            application/json:
              schema:
                type: object
                required:
                  - coding
                properties:
                  coding:
                    type: object
                    required:
                      - completeLangList
                      - genTestLangList
                      - convertLangList
                      - explainLangList
                      - fixErrorLangList
                    properties:
                      completeLangList:
                        type: array
                        items:
                          $ref: '#/components/schemas/Lang'
                        description: 代码补全支持的编程语言列表
                      genTestLangList:
                        type: array
                        items:
                          $ref: '#/components/schemas/Lang'
                        description: 生成测试代码支持的编程语言列表
                      convertLangList:
                        type: array
                        items:
                          $ref: '#/components/schemas/Lang'
                        description: 转换编程语言的支持列表
                      explainLangList:
                        type: array
                        items:
                          $ref: '#/components/schemas/Lang'
                        description: 代码解释支持的编程语言列表
                      fixErrorLangList:
                        type: array
                        items:
                          $ref: '#/components/schemas/Lang'
                        description: 根据错误提示给出修复方案支持的编程语言列表
        '403':
          description: 没有权限
          headers:
            Access-Control-Allow-Origin:
              schema:
                type: string
                default: '*'
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrInfo'
        '500':
          description: 失败
          headers:
            Access-Control-Allow-Origin:
              schema:
                type: string
                default: '*'
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrInfo'
  /api/dev/genToken:
    post:
      tags:
        - dev
      summary: 生成authToken,只在开发模式下有效
      description: 生成authToken,只在开发模式下有效
      operationId: apiDevGenTokenPost
      requestBody:
        required: true
        content:
          application/x-www-form-urlencoded:
            schema:
              type: object
              required:
                - contextValue
                - randomStr
              properties:
                contextValue:
                  type: string
                  description: 上下文信息
                randomStr:
                  type: string
                  description: 随机字符串，加密因子。需要32位长度以上
      responses:
        '200':
          description: 成功
          headers:
            Access-Control-Allow-Origin:
              schema:
                type: string
                default: '*'
          content:
            application/json:
              schema:
                type: object
                required:
                  - token
                properties:
                  token:
                    type: string
                    description: 登录凭证
        '500':
          description: 失败
          headers:
            Access-Control-Allow-Origin:
              schema:
                type: string
                default: '*'
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrInfo'
components:
  parameters:
    AuthToken:
      in: header
      name: X-AuthToken
      schema:
        type: string
      required: true
    DestLang:
      in: query
      name: destLang
      schema:
        type: string
        enum:
          - python
          - c
          - cplusplus
          - java
          - csharp
          - visualbasic
          - javascript
          - sql
          - asm
          - php
          - r
          - go
          - matlab
          - swift
          - delphi
          - ruby
          - perl
          - objc
          - rust
      required: true
    Lang:
      in: path
      name: lang
      schema:
        type: string
        enum:
          - python
          - c
          - cplusplus
          - java
          - csharp
          - visualbasic
          - javascript
          - sql
          - asm
          - php
          - r
          - go
          - matlab
          - swift
          - delphi
          - ruby
          - perl
          - objc
          - rust
      required: true
  schemas:
    ErrInfo:
      type: object
      properties:
        errMsg:
          type: string
          description: 错误信息
    Lang:
      type: string
      description: 编程语言
      enum:
        - python
        - c
        - cplusplus
        - java
        - csharp
        - visualbasic
        - javascript
        - sql
        - asm
        - php
        - r
        - go
        - matlab
        - swift
        - delphi
        - ruby
        - perl
        - objc
        - rust
`;