---
title: Java7_8中的HashMap深挖.md
author: NuyoahCh
date: 2025-08-06 19:15:43
description: Java7_8中的HashMap深挖
tags: [Java]
categories:
 - [Java] ##目录
 - [Java随笔] ##目录
top: false
cover:
---

<h1 id="gni3P">Java7 HashMap</h1>
<font style="color:rgba(0, 0, 0, 0.65);">HashMap 是最简单的，一来我们非常熟悉，二来就是它不支持并发操作，所以源码也非常简单。</font>

![](https://blog-1328413179.cos.ap-beijing.myqcloud.com/images/image.png?q-sign-algorithm=sha1&q-ak=AKIDjyP9s57jHds9ZOkXbe0TX8dL4-CoFPmv9Csn_bhjdN6FZdVXh2D7VdK581NoWaAy&q-sign-time=1754479269;1754482869&q-key-time=1754479269;1754482869&q-header-list=host&q-url-param-list=ci-process&q-signature=b4d75b886cec310e8d1d2c812d22b27ca59e75f2&x-cos-security-token=2Vp3B9gTqPQMVUfvQ5sK2lorIlGF26Cada4ee1b1ad3ba59e8c1dba0fc5450e038QoJzvtssLVYv9uVpFpRmZq34TUwIA_2hVKZAVNWLaM25TomAiFIDtZrC3nYD4Kw23kUsZLd1p1YjvrFy6ue5bBjbd7kcZ0XZM2I-vTurZhnswsRMSrZ18yg2KKwc-QMPXKHFDbRNQCZs3SB5OEqgGI-3VLKBGqmKm-JjmBMUmZRwdE4RILlz7d7MM_T3v68HyW9R6LwstVRobsYETEdY4baD_4zqQrZ5pcDbPooEIYoa0CrhCnzgnPZlRQadCCrsTx2_GToF8Ta6AEtFlCmWg&ci-process=originImage)

<font style="color:rgba(0, 0, 0, 0.65);">大方向上，HashMap 里面是一个</font>**<font style="color:rgba(0, 0, 0, 0.65);">数组</font>**<font style="color:rgba(0, 0, 0, 0.65);">，然后数组中每个元素是一个</font>**<font style="color:rgba(0, 0, 0, 0.65);">单向链表</font>**<font style="color:rgba(0, 0, 0, 0.65);">。</font>

<font style="color:rgba(0, 0, 0, 0.65);">上图中，每个绿色的实体是嵌套类 Entry 的实例，Entry 包含四个属性：key, value, hash 值和用于单向链表的 next。</font>

<font style="color:rgba(0, 0, 0, 0.65);">capacity：当前数组容量，始终保持 2^n，可以扩容，扩容后数组大小为当前的 2 倍。</font>

> <font style="color:rgba(0, 0, 0, 0.65);">2^n 的根本原因还是落实到取模还是位运算上的效率问题</font>
>

<font style="color:rgba(0, 0, 0, 0.65);">loadFactor：负载因子，默认为 0.75。</font>

> <font style="color:rgba(0, 0, 0, 0.65);">经过大数据下的哈希碰撞等相关测试得出的合理数值结论</font>
>

<font style="color:rgba(0, 0, 0, 0.65);">threshold：扩容的阈值，等于 capacity * loadFactor</font>

<h2 id="IdjTI"><font style="color:rgba(0, 0, 0, 0.65);">Put 过程分析</font></h2>
```java
public V put(K key, V value) {
    // 当插入第一个元素的时候，需要先初始化数组大小
    if (table == EMPTY_TABLE) {
        inflateTable(threshold);
    }
    // 如果 key 为 null，感兴趣的可以往里看，最终会将这个 entry 放到 table[0] 中
    if (key == null)
        // Null Key 的操作方法
        return putForNullKey(value);
    // 1. 求 key 的 hash 值
    int hash = hash(key);
    // 2. 找到对应的数组下标
    int i = indexFor(hash, table.length);
    // 3. 遍历一下对应下标处的链表，看是否有重复的 key 已经存在，
    //    如果有，直接覆盖，put 方法返回旧值就结束了
    for (Entry<K,V> e = table[i]; e != null; e = e.next) {
        Object k;
        if (e.hash == hash && ((k = e.key) == key || key.equals(k))) {
            V oldValue = e.value;
            e.value = value;
            e.recordAccess(this);
            return oldValue;
        }
    }

    modCount++;
    // 4. 不存在重复的 key，将此 entry 添加到链表中，细节后面说
    addEntry(hash, key, value, i);
    return null;
}
```

<h2 id="u8gXp">数组初始化</h2>
<font style="color:rgba(0, 0, 0, 0.65);">在第一个元素插入 HashMap 的时候做一次数组的初始化，就是先确定初始的数组大小，并计算数组扩容的阈值。</font>

```java
private void inflateTable(int toSize) {
    // 保证数组大小一定是 2 的 n 次方。
    // 比如这样初始化：new HashMap(20)，那么处理成初始数组大小是 32
    int capacity = roundUpToPowerOf2(toSize);
    // 计算扩容阈值：capacity * loadFactor
    threshold = (int) Math.min(capacity * loadFactor, MAXIMUM_CAPACITY + 1);
    // 算是初始化数组吧
    table = new Entry[capacity];
    initHashSeedAsNeeded(capacity); //ignore
}
```

<h2 id="GgFx4">计算具体数组位置</h2>
<font style="color:rgba(0, 0, 0, 0.65);">这个简单，我们自己也能 YY 一个：使用 key 的 hash 值对数组长度进行取模就可以了。</font>

```java
static int indexFor(int hash, int length) {
    // assert Integer.bitCount(length) == 1 : "length must be a non-zero power of 2";
    // 这个地方就显示出 2 ^ n 的重要性了
    return hash & (length-1);
}
```

<font style="color:rgba(0, 0, 0, 0.65);">这个方法很简单，简单说就是取 hash 值的低 n 位。如在数组长度为 32 的时候，其实取的就是 key 的 hash 值的低 5 位，作为它在数组中的下标位置。</font>

<h2 id="TvU30"><font style="color:rgba(0, 0, 0, 0.7);">添加节点到链表中</font></h2>
<font style="color:rgba(0, 0, 0, 0.65);">找到数组下标后，会先进行 key 判重，如果没有重复，就准备将新值放入到链表的</font>**<font style="color:rgba(0, 0, 0, 0.65);">表头</font>**<font style="color:rgba(0, 0, 0, 0.65);">。</font>

> <font style="color:rgba(0, 0, 0, 0.65);">Java 7 很明显，还是使用头插法的方法进行操作，这种方法在 Java 8 中就进行了相关的改进，开始使用尾插法的方式，而且操作更加简单，同时防止了死循环的现象</font>
>

```java
void addEntry(int hash, K key, V value, int bucketIndex) {
    // 如果当前 HashMap 大小已经达到了阈值，并且新值要插入的数组位置已经有元素了，那么要扩容
    if ((size >= threshold) && (null != table[bucketIndex])) {
        // 扩容，后面会介绍一下
        resize(2 * table.length);
        // 扩容以后，重新计算 hash 值
        hash = (null != key) ? hash(key) : 0;
        // 重新计算扩容后的新的下标
        bucketIndex = indexFor(hash, table.length);
    }
    // 往下看
    createEntry(hash, key, value, bucketIndex);
}
// 这个很简单，其实就是将新值放到链表的表头，然后 size++
void createEntry(int hash, K key, V value, int bucketIndex) {
    Entry<K,V> e = table[bucketIndex];
    table[bucketIndex] = new Entry<>(hash, key, value, e);
    size++;
}
```

<font style="color:rgba(0, 0, 0, 0.65);">这个方法的主要逻辑就是先判断是否需要扩容，需要的话先扩容，然后再将这个新的数据插入到扩容后的数组的相应位置处的链表的表头。</font>

<h2 id="QXbiU"><font style="color:rgba(0, 0, 0, 0.65);">数组扩容</font></h2>
<font style="color:rgba(0, 0, 0, 0.65);">前面我们看到，在插入新值的时候，如果</font>**<font style="color:rgba(0, 0, 0, 0.65);">当前的 size 已经达到了阈值，并且要插入的数组位置上已经有元素</font>**<font style="color:rgba(0, 0, 0, 0.65);">，那么就会触发扩容，扩容后，数组大小为原来的 2 倍。</font>

```java
void resize(int newCapacity) {
    // 本质上还是一个 copy 的过程
    Entry[] oldTable = table;
    int oldCapacity = oldTable.length;
    if (oldCapacity == MAXIMUM_CAPACITY) {
        threshold = Integer.MAX_VALUE;
        return;
    }
    // 新的数组
    Entry[] newTable = new Entry[newCapacity];
    // 将原来数组中的值迁移到新的更大的数组中
    transfer(newTable, initHashSeedAsNeeded(newCapacity));
    table = newTable;
    threshold = (int)Math.min(newCapacity * loadFactor, MAXIMUM_CAPACITY + 1);
}
```

<font style="color:rgba(0, 0, 0, 0.65);">扩容就是用一个新的大数组替换原来的小数组，并将原来数组中的值迁移到新的数组中。</font>

<font style="color:rgba(0, 0, 0, 0.65);">由于是双倍扩容，迁移过程中，会将原来 table[i] 中的链表的所有节点，分拆到新的数组的 newTable[i] 和 newTable[i + oldLength] 位置上。如原来数组长度是 16，那么扩容后，原来 table[0] 处的链表中的所有元素会被分配到新数组中 newTable[0] 和 newTable[16] 这两个位置。代码比较简单，这里就不展开了。</font>

<h2 id="Rqbe2"><font style="color:rgba(0, 0, 0, 0.7);">get 过程分析</font></h2>
<font style="color:rgba(0, 0, 0, 0.65);">相对于 put 过程，get 过程是非常简单的。</font>

1. <font style="color:rgba(0, 0, 0, 0.65);">根据 key 计算 hash 值。</font>
2. <font style="color:rgba(0, 0, 0, 0.65);">找到相应的数组下标：hash & (length - 1)。</font>
3. <font style="color:rgba(0, 0, 0, 0.65);">遍历该数组位置处的链表，直到找到相等(==或equals)的 key。</font>

```java
public V get(Object key) {
    // 之前说过，key 为 null 的话，会被放到 table[0]，所以只要遍历下 table[0] 处的链表就可以了
    if (key == null)
        return getForNullKey();
    // 
    Entry<K,V> entry = getEntry(key);

    return null == entry ? null : entry.getValue();
}
```

<font style="color:rgba(0, 0, 0, 0.65);">getEntry(key):</font>

```java
final Entry<K,V> getEntry(Object key) {
    if (size == 0) {
        return null;
    }

    int hash = (key == null) ? 0 : hash(key);
    // 确定数组下标，然后从头开始遍历链表，直到找到为止
    for (Entry<K,V> e = table[indexFor(hash, table.length)];
         e != null;
         e = e.next) {
        Object k;
        if (e.hash == hash &&
            ((k = e.key) == key || (key != null && key.equals(k))))
            return e;
    }
    return null;
}
```

<h1 id="KzfMY">Java8 HashMap</h1>
<font style="color:rgba(0, 0, 0, 0.65);">Java8 对 HashMap 进行了一些修改，最大的不同就是利用了红黑树，所以其由 </font>**<font style="color:rgba(0, 0, 0, 0.65);">数组+链表+红黑树</font>**<font style="color:rgba(0, 0, 0, 0.65);"> 组成。</font>

<font style="color:rgba(0, 0, 0, 0.65);">根据 Java7 HashMap 的介绍，我们知道，查找的时候，根据 hash 值我们能够快速定位到数组的具体下标，但是之后的话，需要顺着链表一个个比较下去才能找到我们需要的，时间复杂度取决于链表的长度，为 O(n)。</font>

<font style="color:rgba(0, 0, 0, 0.65);">为了降低这部分的开销，在 Java8 中，当链表中的元素达到了 8 个时，会将链表转换为红黑树，在这些位置进行查找的时候可以降低时间复杂度为 O(logN)。</font>

> <font style="color:rgba(0, 0, 0, 0.65);">在我看来，Java 8 中对 HashMap 的改进并没有特别大，更多的是优化、完善。</font>
>

<font style="color:rgba(0, 0, 0, 0.65);">来一张图简单示意一下吧：</font>

![](https://blog-1328413179.cos.ap-beijing.myqcloud.com/images/image%20%281%29.png?q-sign-algorithm=sha1&q-ak=AKID-pp2eb7ol2pdReiVZ2oEZcc98OQZ267KTX3o5b5HQ3OZ3M61peX0UKBluuds1PJ3&q-sign-time=1754479294;1754482894&q-key-time=1754479294;1754482894&q-header-list=host&q-url-param-list=ci-process&q-signature=19a07a4d5f21c163e83e839c5369c7b4f3d906ca&x-cos-security-token=2Vp3B9gTqPQMVUfvQ5sK2lorIlGF26Ca788637f69ded3e66e555a9d12fc6d4f88QoJzvtssLVYv9uVpFpRmRJ0WEmKEq-eXkE4MozLDz2wgY6qMdSDwPsXoiZhWr-p9VWWwkXaJB8yJ0paipEtyd9xT9Duhr3hJvbg9DAiRX6NJYBdLzKSv-RmwXqfmdSqpFpY0tLMG-4lpPa1-PCEo27ihK3bLI52JEM4tYJMqZQInWnOD1xrsiBX5RECOGgTXHf0wgw00pvwjDcDKV9Hpx7G5EnuqJxEvgHlGyrjNK54KkW9UpSN9myHw2cS5dJPIW8QW_CeP3xTjkbpua61oA&ci-process=originImage)

<font style="color:rgba(0, 0, 0, 0.65);">Java7 中使用 Entry 来代表每个 HashMap 中的数据节点，Java8 中使用 </font>**<font style="color:rgba(0, 0, 0, 0.65);">Node</font>**<font style="color:rgba(0, 0, 0, 0.65);">，基本没有区别，都是 key，value，hash 和 next 这四个属性，不过，Node 只能用于链表的情况，红黑树的情况需要使用 </font>**<font style="color:rgba(0, 0, 0, 0.65);">TreeNode</font>**<font style="color:rgba(0, 0, 0, 0.65);">。</font>

<font style="color:rgba(0, 0, 0, 0.65);">我们根据数组元素中，第一个节点数据类型是 Node 还是 TreeNode 来判断该位置下是链表还是红黑树的。</font>

<h2 id="P9YLp"><font style="color:rgba(0, 0, 0, 0.65);">Put 过程分析</font></h2>
```java
public V put(K key, V value) {
    return putVal(hash(key), key, value, false, true);
}

// 第四个参数 onlyIfAbsent 如果是 true，那么只有在不存在该 key 时才会进行 put 操作
// 第五个参数 evict 我们这里不关心
final V putVal(int hash, K key, V value, boolean onlyIfAbsent,
               boolean evict) {
    Node<K,V>[] tab; Node<K,V> p; int n, i;
    // 第一次 put 值的时候，会触发下面的 resize()，类似 java7 的第一次 put 也要初始化数组长度
    // 第一次 resize 和后续的扩容有些不一样，因为这次是数组从 null 初始化到默认的 16 或自定义的初始容量
    if ((tab = table) == null || (n = tab.length) == 0)
        n = (tab = resize()).length;
    // 找到具体的数组下标，如果此位置没有值，那么直接初始化一下 Node 并放置在这个位置就可以了
    if ((p = tab[i = (n - 1) & hash]) == null)
        tab[i] = newNode(hash, key, value, null);
  
    else {// 数组该位置有数据
        Node<K,V> e; K k;
        // 首先，判断该位置的第一个数据和我们要插入的数据，key 是不是"相等"，如果是，取出这个节点
        if (p.hash == hash &&
            ((k = p.key) == key || (key != null && key.equals(k))))
            e = p;
        // 如果该节点是代表红黑树的节点，调用红黑树的插值方法，本文不展开说红黑树
        else if (p instanceof TreeNode)
            e = ((TreeNode<K,V>)p).putTreeVal(this, tab, hash, key, value);
        else {
            // 到这里，说明数组该位置上是一个链表
            for (int binCount = 0; ; ++binCount) {
                // 插入到链表的最后面(Java7 是插入到链表的最前面)
                if ((e = p.next) == null) {
                    p.next = newNode(hash, key, value, null);
                    // TREEIFY_THRESHOLD 为 8，所以，如果新插入的值是链表中的第 8 个
                    // 会触发下面的 treeifyBin，也就是将链表转换为红黑树
                    if (binCount >= TREEIFY_THRESHOLD - 1) // -1 for 1st
                        treeifyBin(tab, hash);
                    break;
                }
                // 如果在该链表中找到了"相等"的 key(== 或 equals)
                if (e.hash == hash &&
                    ((k = e.key) == key || (key != null && key.equals(k))))
                    // 此时 break，那么 e 为链表中[与要插入的新值的 key "相等"]的 node
                    break;
                p = e;
            }
        }
        // e!=null 说明存在旧值的key与要插入的key"相等"
        // 对于我们分析的put操作，下面这个 if 其实就是进行 "值覆盖"，然后返回旧值
        if (e != null) {
            V oldValue = e.value;
            if (!onlyIfAbsent || oldValue == null)
                e.value = value;
            afterNodeAccess(e);
            return oldValue;
        }
    }
    // 判断并发的处理
    ++modCount;
    // 如果 HashMap 由于新插入这个值导致 size 已经超过了阈值，需要进行扩容
    if (++size > threshold)
        resize();
    afterNodeInsertion(evict);
    return null;
}
```

<font style="color:rgba(0, 0, 0, 0.65);">和 Java7 稍微有点不一样的地方就是，Java7 是先扩容后插入新值的，Java8 先插值再扩容，不过这个不重要。</font>

<h2 id="dEhVN">数组扩容</h2>
<font style="color:rgba(0, 0, 0, 0.65);">resize() 方法用于</font>**<font style="color:rgba(0, 0, 0, 0.65);">初始化数组</font>**<font style="color:rgba(0, 0, 0, 0.65);">或</font>**<font style="color:rgba(0, 0, 0, 0.65);">数组扩容</font>**<font style="color:rgba(0, 0, 0, 0.65);">，每次扩容后，容量为原来的 2 倍，并进行数据迁移。</font>

```java
final Node<K,V>[] resize() {
    Node<K,V>[] oldTab = table;
    int oldCap = (oldTab == null) ? 0 : oldTab.length;
    int oldThr = threshold;
    int newCap, newThr = 0;
    if (oldCap > 0) { // 对应数组扩容
        if (oldCap >= MAXIMUM_CAPACITY) {
            threshold = Integer.MAX_VALUE;
            return oldTab;
        }
        // 将数组大小扩大一倍
        else if ((newCap = oldCap << 1) < MAXIMUM_CAPACITY &&
                 oldCap >= DEFAULT_INITIAL_CAPACITY)
            // 将阈值扩大一倍
            newThr = oldThr << 1; // double threshold
    }
    else if (oldThr > 0) // 对应使用 new HashMap(int initialCapacity) 初始化后，第一次 put 的时候
        newCap = oldThr;
    else {// 对应使用 new HashMap() 初始化后，第一次 put 的时候
        newCap = DEFAULT_INITIAL_CAPACITY;
        newThr = (int)(DEFAULT_LOAD_FACTOR * DEFAULT_INITIAL_CAPACITY);
    }
    
    if (newThr == 0) {
        float ft = (float)newCap * loadFactor;
        newThr = (newCap < MAXIMUM_CAPACITY && ft < (float)MAXIMUM_CAPACITY ?
                  (int)ft : Integer.MAX_VALUE);
    }
    threshold = newThr;

    // 用新的数组大小初始化新的数组
    Node<K,V>[] newTab = (Node<K,V>[])new Node[newCap];
    table = newTab; // 如果是初始化数组，到这里就结束了，返回 newTab 即可
    
    if (oldTab != null) {
        // 开始遍历原数组，进行数据迁移。
        for (int j = 0; j < oldCap; ++j) {
            Node<K,V> e;
            if ((e = oldTab[j]) != null) {
                oldTab[j] = null;
                // 如果该数组位置上只有单个元素，那就简单了，简单迁移这个元素就可以了
                if (e.next == null)
                    newTab[e.hash & (newCap - 1)] = e;
                // 如果是红黑树，具体我们就不展开了
                else if (e instanceof TreeNode)
                    ((TreeNode<K,V>)e).split(this, newTab, j, oldCap);
                else { 
                    // 这块是处理链表的情况，
                    // 需要将此链表拆成两个链表，放到新的数组中，并且保留原来的先后顺序
                    // loHead、loTail 对应一条链表，hiHead、hiTail 对应另一条链表，代码还是比较简单的
                    Node<K,V> loHead = null, loTail = null;
                    Node<K,V> hiHead = null, hiTail = null;
                    Node<K,V> next;
                    do {
                        next = e.next;
                        if ((e.hash & oldCap) == 0) {
                            if (loTail == null)
                                loHead = e;
                            else
                                loTail.next = e;
                            loTail = e;
                        }
                        else {
                            if (hiTail == null)
                                hiHead = e;
                            else
                                hiTail.next = e;
                            hiTail = e;
                        }
                    } while ((e = next) != null);
                    if (loTail != null) {
                        loTail.next = null;
                        // 第一条链表
                        newTab[j] = loHead;
                    }
                    if (hiTail != null) {
                        hiTail.next = null;
                        // 第二条链表的新的位置是 j + oldCap，这个很好理解
                        newTab[j + oldCap] = hiHead;
                    }
                }
            }
        }
    }
    return newTab;
}
```

<h2 id="pFKME"><font style="color:rgba(0, 0, 0, 0.7);">get 过程分析</font></h2>
<font style="color:rgba(0, 0, 0, 0.65);">相对于 put 来说，get 真的太简单了。</font>

1. <font style="color:rgba(0, 0, 0, 0.65);">计算 key 的 hash 值，根据 hash 值找到对应数组下标: hash & (length-1)</font>
2. <font style="color:rgba(0, 0, 0, 0.65);">判断数组该位置处的元素是否刚好就是我们要找的，如果不是，走第三步</font>
3. <font style="color:rgba(0, 0, 0, 0.65);">判断该元素类型是否是 TreeNode，如果是，用红黑树的方法取数据，如果不是，走第四步</font>
4. <font style="color:rgba(0, 0, 0, 0.65);">遍历链表，直到找到相等(==或equals)的 key</font>

```java
public V get(Object key) {
    Node<K,V> e;
    return (e = getNode(hash(key), key)) == null ? null : e.value;
}
```

```java
final Node<K,V> getNode(int hash, Object key) {
    Node<K,V>[] tab; Node<K,V> first, e; int n; K k;
    if ((tab = table) != null && (n = tab.length) > 0 &&
        (first = tab[(n - 1) & hash]) != null) {
        // 判断第一个节点是不是就是需要的
        if (first.hash == hash && // always check first node
            ((k = first.key) == key || (key != null && key.equals(k))))
            return first;
        if ((e = first.next) != null) {
            // 判断是否是红黑树
            if (first instanceof TreeNode)
                return ((TreeNode<K,V>)first).getTreeNode(hash, key);

            // 链表遍历
            do {
                if (e.hash == hash &&
                    ((k = e.key) == key || (key != null && key.equals(k))))
                    return e;
            } while ((e = e.next) != null);
        }
    }
    return null;
}
```
