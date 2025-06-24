---
title: 実装しながら学ぶ！HTTPサーバーの基本 / 技育CAMPアカデミア
author: Takuma Kobayashi @ Finatext
header: 実装しながら学ぶ！HTTPサーバーの基本 / 技育CAMPアカデミア
slide: true
marp: true
theme: meta
---

# HTTP サーバー実装の基本

Takuma Kobayashi @(@takuma5884rb)
株式会社 Finatext

---

## 目次

<!-- paginate: true -->

- 自己紹介
- 今日学んで欲しいこと
- リクエストを受け取る
- 異常系を考える
- 認証
- まとめ

---

## 自己紹介

- 小林拓磨
  - X: @takuma58884rbb
- Software Engineer at Finatext
  - １年目からシステムの詳細設計・実装・運用を経験
  - 主な技術スタック
    - Go, AWS, Terraform
- 2024 Japan AWS Jr.Champions
- 趣味は料理・マラソン

![bg right:45%](./TAKUMA.jpeg)

---

## 今日学んで欲しいこと

---

## 今日学んで欲しいこと

1. **関数の責務を捉える**

   - 関心事の分離による保守性の向上
   - 責任範囲の明確化

2. **適切なリクエスト処理とバリデーション**

   - 不正なリクエストからサーバーを守る
   - エラーハンドリングの重要性

3. **認証の基本と実装方法**
   - セキュアな認証の実現
   - トークンベース認証のメリット・デメリット

---

## レイヤードアーキテクチャ

---

**クリーンアーキテクチャとは誰々が提唱した概念で〜**

![bg 70% blur:3px opacity:.4](./CleanArchitecture.jpg)

**...という話は今日はしません！**

> > > https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

---

## 宣伝

---

![bg 80%](./summer_internship_summer_internship_engineer.png)

---

## ご清聴ありがとうございました！

質問やフィードバックがあればお気軽にどうぞ！
X の DM でも受け付けています！
![bg right:35% width:450px](./x.png)
