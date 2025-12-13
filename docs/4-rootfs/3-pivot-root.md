# 4-3. pivot_rootを用いたルート移動

前節で扱ったchrootは、**比較的弱い**機構です。  
もしコンテナ内のユーザーが**chrootする権限**を持っていた場合、元のルートディレクトリに移動できてしまいます。

詳しくはこちら:  
<https://container-security.dev/namespace/chroot-and-pivot_root#chroot-%E3%81%AE%E5%95%8F%E9%A1%8C%E7%82%B9>

少々**扱いは難しい**ですが**より強固**なルートディレクトリの移動方法として、**pivot_root**があります。  
さあ、pivot_rootで**より強固なセキュリティ**を手に入れましょう。

## 【前提】この節で扱うsyscall

### [pivot_root](https://pkg.go.dev/golang.org/x/sys/unix#PivotRoot)

ルートファイルシステムの**マウントを入れ替え**ます。  
より正確には、**今までのルートファイルシステム**を`putold`のディレクトリにマウントし、`newroot`のディレクトリを**新しいルートファイルシステム**として`/`にマウントします。

[4-1](/4-rootfs/1-mount-rootfs.md)で書いた通り、ルートファイルシステムも**マウントによって構築**されています。  
このマウントを**入れ替えて**しまえば、抜ける手段はないだろうという算段です。  
OS側で登録されているルートフォルダを変更するだけのchrootと違い、**比較的安全**なルート移動方法です。

```go
func PivotRoot(newroot string, putold string) (err error)
```

## ルートを移動させる

pivot_rootには**いくつかの制約**があります。

> `new_root`および`put_old`には以下の制限がある:
>
> - ディレクトリでなければならない。
> - `new_root`と`put_old`は現在の`root`と同じファイルシステムにあってはならない。
> - `put_old`は`new_root`以下になければならない。

さて、これらの制約をクリアするにはどうすれば良いでしょうか？  
chrootよりもはるかに複雑なマウント処理が必要になりますが、頑張ってみましょう。
