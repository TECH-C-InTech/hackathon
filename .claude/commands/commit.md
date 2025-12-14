---
name: commit
description: 変更をコミット（Conventional Commits準拠）
---

現在の変更を確認して、Conventional Commitsに従ったコミットメッセージを作成してコミットして。

## コミット粒度
- 1コミット = 1つの論理的な変更
- リファクタと機能追加は分離
- テストコードは本体と同じコミットでOK
- フォーマット変更のみなら別コミット

## コミットメッセージ形式
- **1行のみ**（bodyは不要）
- 形式: `type: 変更内容`
- bodyが必要なら粒度が大きすぎる。分割しろ

必ず以下を追加:
```
Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
```

