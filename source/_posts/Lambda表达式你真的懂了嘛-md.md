---
title: Lambda表达式你真的懂了嘛.md
author: NuyoahCh
date: 2025-07-27 22:32:57
description: Lambda表达式你真的懂了嘛
tags: [Java]
categories:
 - [Java] ##目录
 - [Java随笔] ##目录
top: true
# cover:
---

<h1 id="hj3Um"><font style="color:rgb(37, 41, 51);">一、Lambda 表达式介绍</font></h1>

<font style="color:rgb(37, 41, 51);">Lambda 表达式是 Java 8 中添加的功能。引入 Lambda 表达式的主要目的是为了让 Java 支持函数式编程。 Lambda 表达式是一个可以在不属于任何类的情况下创建的函数，并且可以像对象一样被传递和执行。</font>

**<font style="color:rgb(37, 41, 51);">Java lambda 表达式用于实现简单的单方法接口，与 Java Streams API 配合进行函数式编程</font>**<font style="color:rgb(37, 41, 51);">。</font>

<font style="color:rgb(37, 41, 51);">以一个例子引入 lambda 表达式</font>

```java
/**
 * @author 小王 Coding
 */
public class App {
    public static void main(String[] args) {
        Runnable runnable = new Runnable() {
            @Override
            public void run() {
                System.out.println("这是匿名内部类的实现");
            }
        };

        Runnable runnable1 = () -> System.out.println("这是Lambda表达式的实现");

        runnable.run();
        runnable1.run();
    }
}
```

![](https://cdn.nlark.com/yuque/0/2025/png/45054063/1753620668358-bbedc240-e8db-4981-bf4d-8da7518e4caa.png)

```java
// 重写 Runnable 接口 run 方法进行打印
@Override
public void run() {
    System.out.println("这是匿名内部类的实现");
}
```

相比较匿名内部类的重写接口抽象方法的实现，Lambda 表达式展现的更加优雅

![](https://cdn.nlark.com/yuque/0/2025/png/45054063/1753625562056-bf173b28-b3df-45de-a76d-3e944127ed14.png)

> 那么到底是普通实现的方式更好，还是 ALL IN Lambda 表达式呢？？？
>
> 别急，现在开始 Lambda 表达式的详细剖析～

<h2 id="zoxqA">Lambda 表达式</h2>

Lambda是简洁的标识可传递匿名函数的一种方式。“互动”事件驱动下，最终面向对象编程和函数式编程结合才是趋势。 java中，一段代码的传递并不容易。因为JAVA是面向对象的语言，如果要传递一段代码，必须先构建类，再生成对应的对象来传递所要的代码。在之前，JAVA的设计者都抗拒加入这一特性，虽然JAVA现有的特性也能通过类和对象实现类似的API但是这样复杂且不易于使用。在后期，问题早已不是JAVA是不是要变成一门使用函数式编程的语言，而是如何实现这种改变。在java8之前已经有了多年的实验，然后JAVA8来了。



<h2 id="h58WK">Lambda 特性</h2>

1. 匿名：lambda表达式不像面向对象的方法一样，有确定的名称。
2. 函数：虽然lambda不是对象的方法，属于某个特定的类。但是lambda表达式一样的有参数列表、函数主体 返回类型和异常声明
3. 传递：lambda表达式可以作为参数传递
4. 简洁：无需像匿名类一样有固定模板的代码，lambda写得少而想得多
5. JAVA8中 可以为接口增加静态方法、可以为类增加默认方法



<h2 id="qg1Af">Lambda 表达式介绍</h2>
<h3 id="I5U0E">Lambda 表达式结构</h3>

![](https://cdn.nlark.com/yuque/0/2025/png/45054063/1753621748707-d3c27f00-ed9b-4cc6-a867-7b4fb367d64b.png)

**基本语法**

+ <font style="color:rgb(77, 77, 77);">(参数…）-> 表达式 隐式返回表达式结果</font>
+ <font style="color:rgb(77, 77, 77);">(参数…)->{执行语句} 可用return语句 显示返回执行结果</font>
+ <font style="color:rgb(77, 77, 77);">函数式接口不允许抛出受检异常</font>
+ <font style="color:rgb(77, 77, 77);">注意：当参数只有一个时，也可以去掉参数的括号。原因是java编译器的自动类型推断</font>



<h3 id="AG2TG">常见的 Lambda 表达式</h3>

```java
//1、单个参数
(String s)->s.length()

//2、单个对象
(Apple a)->a.getWeight()>150

//3、多参数,多语句
(int a,int b)->{
	System.out.println(a);
	System.out.println(b);
}

//4、空参数,返回int值42
()->42

//5、多对象参数
(Applea1,Applea2)->a1.getWeight().compareTo(a2.getWeight())
```

上述的表达式整体上比较简单，这里主要是理解和熟悉 Lambda 表达式的结构



<h3 id="DY1yV">Lambda 表达式阅读技巧</h3>

> 接下来我们使用这个比较复杂的例子，对 Lambda 表达式进行深度剖析！

先看一下这个例子：

```java
List<A> list = aList.filter(a -> a.getId() > 10).collect(Colletors.toList);
```

其中`filter`方法里用到的`a -> a.getId() > 10`就是一个 Lambda 表达式。

这里在详细解释一下这段 lambda 表达式的意思：

+ **aList.stream()**: 将`aList`转换为一个Stream，以便使用Stream API的操作。
+ **filter(a -> a.getId() > 10)**: 使用Lambda表达式作为过滤条件，筛选出`id`大于10的元素。
+ **collect(Collectors.toList())**: 将过滤后的Stream收集到一个List中。
+ 最终，`list`将包含`aList`中所有满足`id > 10`的元素。

这么读起来是不是也不是那么复杂了呢～

> 这里也分享一个阅读 Lambda 表达式的小技巧
>
> **把 Lambda 表达式“拆成方法看”：****<font style="color:#0e0e0e;">看到 Lambda，就想象它背后其实是一个匿名内部类 + 实现接口的方法体</font>**

**<font style="color:#0e0e0e;"></font>**

<h2 id="bneBY"><font style="color:#0e0e0e;">Lambda 表达式注意事项</font></h2>
<h3 id="VxT8R">类型检查</h3>

1. Lambda的类型由上下文推断而来
2. 同样的lambda表达式，不同的函数式接口，只要方法的签名一致，同样的表达式可用于不同的函数是接口。
3. 只有函数式接口的实现，能承载lambda表达式
4. Objecto=()-{System.out.print("HellowWorld")} 这是不合法的 因为Object不是一个函数式接口



<h3 id="YGZGU">类型推断</h3>

1. <font style="color:rgb(77, 77, 77);">Lambda表达式可以省略参数的类型，java编译器能自动推断</font>
2. <font style="color:rgb(77, 77, 77);">当lambda只有一个参数需要推断类型时，参数两边的括号可以省略</font>

```java
List<Apple> c=filter(inventory,a->"green".equals(a.getColor()));
Comparator<Apple> c=(a1,a2)->a1.getWeight.compareTo(a2.getWeight());
```



<h3 id="w7FnX"><font style="color:rgb(79, 79, 79);">变量作用域</font></h3>

JAVA8之前 内部类只允许访问final修饰的变量，现在使用lambda表达式，一个内部类可以访问任何有效的final局部变量-任何值不会发生变化的变量

1. java限制了 lambda表达式访问的自由变量，值是不可更改的，因为这会导致出现无法预料的并发问题。
2. java编译器的限制是有限的，只对局部变量有效，如果使用静态变量，或者示例变量，编译器不会提示任何错误。这样仍然是不安全的。
3. 可以用数组 `int[] counter =new int[1]; button.setOnAaction(event->counter[0]++);`仍然 可以让lambda对局部变量进行重新赋值。
4. lambda表达式的方法体，与被嵌套的代码块具有同样的作用域，所有适用同样的命名冲突和变量屏蔽规则。



<h3 id="zAcGr">方法引用</h3>

对于已有的方法，如果希望作为lambda表达式来使用，可以直接使用方法引用

三种方法引用的情况

1. 对象：：实例方法
2. 类：：静态方法
3. 类：：实例方法

在第一种和第二种方法引用种，方法的引用等于提供方法参数的lambda表达式

例如：

`System.out::println() 等同于 System.out.print(x)`

`Math::pow 等同于 （x,y）->Math.pow(x,y)`

对于第三种，则相当于第一个参数成为执行方法的对象

例如：`String::compareToIngnoreCase 等同于（x,y） x.compareIngoreCase（Y）；`



<h3 id="e6sMY">构造器引用</h3>

对于构造器引用，相当于根据构造器的方法的参数，生成一个构造的对象的一个lambda表达式

例如：`StringBuilder::new` 可以表示为 `(Stiring s)->new StringBuilder(s);` 具体引用哪个构造器，编译器会根据上下文推断使用符合参数的构造器。



<h1 id="EXPJ4">二、<font style="color:rgb(37, 41, 51);">Lambda 表达式和函数式接口</font></h1>

lambda 表达式便于实现只拥有单一方法的接口，同样在 Java 里匿名类也用于快速实现接口，只不过 lambda 相较于匿名类更方便些，在书写的时候连创建类的步骤也免去了，更适合用在函数式编程。

举个例子来说，函数式编程经常用在实现事件 Listener 的时候 。 在 Java 中的事件侦听器通常被定义为具有单个方法的 Java 接口。下面是一个 Listener 接口示例：

```java
public interface StateChangeListener {

    public void onStateChange(State oldState, State newState);

}
```

上面这个 Java 接口定义了一个只要被监听对象的状态发生变化，就会调用的 onStateChange 方法（这里不用管监听的是什么，举例而已）。 在 Java 8 版本以前，监听事件变更的程序必须实现此接口才能侦听状态更改。

比如说，有一个名为 StateOwner 的类，它可以注册状态的事件侦听器。

```java
public class StateOwner {

    public void addStateListener(StateChangeListener listener) { ... }

}

```

<font style="color:rgb(37, 41, 51);">我们可以使用匿名类实现 StateChangeListener 接口，然后为 StateOwner 实例添加侦听器。</font>

```java
StateOwner stateOwner = new StateOwner();

stateOwner.addStateListener(new StateChangeListener() {

    public void onStateChange(State oldState, State newState) {
        // do something with the old and new state.
        System.out.println("State changed")
    }
});
```

<font style="color:rgb(37, 41, 51);">在 Java 8 引入Lambda 表达式后，我们可以用 Lambda 表达式实现 StateChangeListener 接口会更加方便。现在，把上面例子接口的匿名类实现改为 Lambda 实现，程序会变成这样：</font>

```java
StateOwner stateOwner = new StateOwner();

stateOwner.addStateListener(
    (oldState, newState) -> System.out.println("State changed")
);
```

显而易见，这个变化还是比较明显的

<font style="color:rgb(37, 41, 51);">在这里，我们使用的 Lambda 表达式是：</font>

```java
(oldState, newState) -> System.out.println("State changed")
```

这个 lambda 表达式与 StateChangeListener 接口的 onStateChange() 方法的参数列表和返回值类型相匹配。**如果一个 lambda 表达式匹配单方法接口中方法的参数列表和返回值（比如本例中的 StateChangeListener 接口的 onStateChange 方法），则 lambda 表达式将转换为拥有相同方法签名的接口实现。** 这句话听着有点绕，下面详细解释一下 Lambda 表达式和接口匹配的详细规则。

<h2 id="RCOFE"> 匹配Lambda 与接口的规则</h2>

上面例子里使用的 StateChangeListener 接口有一个特点，其只有一个未实现的抽象方法，在 Java 里这样的接口也叫做函数式接口 （Functional Interface）。将 Java lambda 表达式与接口匹配需要满足一下三个规则：

+ 接口是否只有一个抽象（未实现）方法，即是一个函数式接口？
+ lambda 表达式的参数是否与抽象方法的参数匹配？
+ lambda 表达式的返回类型是否与单个方法的返回类型匹配？

如果能满足这三个条件，那么给定的 lambda 表达式就能与接口成功匹配类型。

<h2 id="rIk6b">函数式接口</h2>

**只有一个抽象方法的接口被称为函数是式接口**，从 Java 8 开始，Java 接口中可以包含默认方法和静态方法。默认方法和静态方法都有直接在接口声明中定义的实现。这意味着，Java lambda 表达式可以实现拥有多个方法的接口——只要接口中只有一个未实现的抽象方法就行。

所以在文章一开头说lambda 用于实现单方法接口，是为了让大家更好的理解，真实的情况是只要接口中只存在一个抽象方法，那么这个接口就能用 lambda 实现。

换句话说，即使接口包含默认方法和静态方法，只要接口只包含一个未实现的抽象方法，它就是函数式接口。比如下面这个接口：

```java
import java.io.IOException;
import java.io.OutputStream;

public interface MyInterface {

    void printIt(String text);

    default public void printUtf8To(String text, OutputStream outputStream){
        try {
            outputStream.write(text.getBytes("UTF-8"));
        } catch (IOException e) {
            throw new RuntimeException("Error writing String as UTF-8 to OutputStream", e);
        }
    }

    static void printItToSystemOut(String text){
        System.out.println(text);
    }
}

```

<font style="color:rgb(37, 41, 51);">即使这个接口包含 3 个方法，它也可以通过 lambda 表达式实现，因为接口中只有一个抽象方法 printIt没有被实现。</font>

```java
MyInterface myInterface = (String text) -> {
    System.out.print(text);
};
```



<h2 id="3a8f7fbd">把方法引用作为 Lambda</h2>

如过编写的 lambda 表达式所做的只是使用传递给 Lambda 的参数调用另一个方法，那么 Java里为 Lambda 实现提供了一种更简短的形式来表达方法调用。比如说，下面是一个函数式数接口：

```java
public interface MyPrinter{
    public void print(String s);
}
```

接下来我们用 Lambda 表达式实现这个 MyPrinter 接口

```java
MyPrinter myPrinter = (s) -> { System.out.println(s); };
```

因为 Lambda 的参数只有一个，方法体也只包含一行，所以可以简写成

```java
MyPrinter myPrinter = s ->  System.out.println(s);
```

又因为 Lambda 方法体内所做的只是将字符串参数转发给 System.out.println() 方法，因此我们可以将上面的 Lambda 声明替换为方法引用。

```java
MyPrinter myPrinter = System.out::println;
```

**注意双冒号 :: 向 Java 的编译器指明这是一个方法的引用。引用的方法是双冒号之后的方法。而拥有引用方法的类或对象则位于双冒号之前。**

我们可以引用以下类型的方法：

+ 静态方法
+ 参数对象的实例方法
+ 实例方法
+ 类的构造方法



<h3 id="bb493c23">引用类的静态方法</h3>

最容易引用的方法是静态方法，比如有这么一个函数式接口和类

```java
public interface Finder {
    public int find(String s1, String s2);
}

public class MyClass{
    public static int doFind(String s1, String s2){
        return s1.lastIndexOf(s2);
    }
}
```

如果我们创建 Lambda 去调用 MyClass 的静态方法 doFind

```java
Finder finder = (s1, s2) -> MyClass.doFind(s1, s2);
```

所以我们可以使用 Lambda 直接引用 Myclass 的 doFind 方法。

```java
Finder finder = MyClass::doFind;
```



<h3 id="89047f28">引用参数的方法</h3>

接下来，如果我们在 Lambda 直接转发调用的方法是来自参数的方法

```java
public interface Finder {
    public int find(String s1, String s2);
}

Finder finder = (s1, s2) -> s1.indexOf(s2);
```

依然可以通过 Lambda 直接引用

```java
Finder finder = String::indexOf;
```

这个与上面完全形态的 Lambda 在功能上完全一样，不过要注意简版 Lambda 是如何引用单个方法的。 Java 编译器会尝试将引用的方法与第一个参数的类型匹配，使用第二个参数类型作为引用方法的参数。



<h3 id="27a9a8f6">引用实例方法</h3>

我们还也可以从 Lambda 定义中引用实例方法。首先，设想有如下接口

```java
public interface Deserializer {
    public int deserialize(String v1);
}
```

该接口表示一个能够将字符串“反序列化”为 int 的组件。现在有一个 StringConvert 类

```java
public class StringConverter {
    public int convertToInt(String v1){
        return Integer.valueOf(v1);
    }
}
```

StringConvert 类 的 convertToInt() 方法与 Deserializer 接口的 deserialize() 方法具有相同的签名。因此，我们可以创建 StringConverter 的实例并从 Lambda 表达式中引用其 convertToInt() 方法，如下所示：

```java
StringConverter stringConverter = new StringConverter();

Deserializer des = stringConverter::convertToInt;
// 等同于 Deserializer des = (value) -> stringConverter.convertToInt(value)
```

上面第二行代码创建的 Lambda 表达式引用了在第一行创建的 StringConverter 实例的 convertToInt 方法。



<h3 id="34ec6e1d">引用构造方法</h3>

最后如果 Lambda 的作用是调用一个类的构造方法，那么可以通过 Lambda 直接引用类的构造方法。在 Lambda 引用类构造方法的形式如下：

```java
ClassName::new
```

那么如何将构造方法用作 lambda 表达式呢，假设我们有这样一个函数式接口

```java
public interface Factory {
    public String create(char[] val);
}
```

Factory 接口的 create() 方法与 String 类中的其中一个构造方法的签名相匹配（String 类有多个重载版本的构造方法）。因此，String类的该构造方法也可以用作 Lambda 表达式。

```java
Factory factory = String::new;
// 等同于 Factory factory (chars) -> String.new(chars);
```



<h2 id="AAm1U"><font style="color:rgb(79, 79, 79);">常见的Lambda和已有的实现</font></h2>

![](https://cdn.nlark.com/yuque/0/2025/png/45054063/1753624375356-a14b65f2-7832-4f0f-ad1d-9b0b5aeff998.png)



<h1 id="bNfWK">三、Lambda 表达式性能问题</h1>

**<font style="color:rgb(51, 51, 51);">有人说“Lambda 能让 Java 程序慢 30 倍”，你怎么看？</font>**

<h2 id="Ly90i"><font style="color:rgb(51, 51, 51);">基准测试表明立场</font></h2>

<font style="color:rgb(51, 51, 51);">为了让你清楚地了解这个背景，请参考下面的代码片段。在实际运行中，基于 Lambda/Stream 的版本（lambdaMaxInteger），比传统的 for-each 版本（forEachLoopMaxInteger）慢很多。</font>

```java
// 一个大的ArrayList，内部是随机的整形数据
volatile List<Integer> integers = …
 
// 基准测试1
public int forEachLoopMaxInteger() {
   int max = Integer.MIN_VALUE;
   for (Integer n : integers) {
    max = Integer.max(max, n);
   }
   return max;
}
 
// 基准测试2
public int lambdaMaxInteger() {
   return integers.stream().reduce(Integer.MIN_VALUE, (a, b) -> Integer.max(a, b));
}
```

<font style="color:rgb(51, 51, 51);">第一，基准测试是一个非常有效的通用手段，让我们以直观、量化的方式，判断程序在特定条件下的性能表现。</font>

<font style="color:rgb(51, 51, 51);">第二，基准测试必须明确定义自身的范围和目标，否则很有可能产生误导的结果。前面代码片段本身的逻辑就有瑕疵，更多的开销是源于自动装箱、拆箱（auto-boxing/unboxing），而不是源自 Lambda 和 Stream，所以得出的初始结论是没有说服力的。</font>

<font style="color:rgb(51, 51, 51);">第三，虽然 Lambda/Stream 为 Java 提供了强大的函数式编程能力，但是也需要正视其局限性：</font>

+ <font style="color:rgb(51, 51, 51);">一般来说，我们可以认为 Lambda/Stream 提供了与传统方式接近对等的性能，但是如果对于性能非常敏感，就不能完全忽视它在特定场景的性能差异了，例如：初始化的开销。 Lambda 并不算是语法糖，而是一种新的工作机制，在首次调用时，JVM 需要为其构建CallSite实例。这意味着，如果 Java 应用启动过程引入了很多 Lambda 语句，会导致启动过程变慢。其实现特点决定了 JVM 对它的优化可能与传统方式存在差异。</font>
+ <font style="color:rgb(51, 51, 51);">增加了程序诊断等方面的复杂性，程序栈要复杂很多，Fluent 风格本身也不算是对于调试非常友好的结构，并且在可检查异常的处理方面也存在着局限性等。</font>



<h2 id="dPig1">个人思考揭开谜团</h2>
<h3 id="a98b892e">正常使用场景下，Lambda 并不会导致明显性能问题</h3>

<font style="color:rgb(51, 51, 51);">在大多数日常使用场景中，比如 forEach、map、filter 等集合操作，</font>**Lambda 的性能几乎等同于匿名内部类，甚至更好（因为 JVM 有优化）**<font style="color:rgb(51, 51, 51);">。</font>

<font style="color:rgb(51, 51, 51);">所以如果你写的是：</font>

```java
list.forEach(item -> System.out.println(item));
```

<font style="color:rgb(51, 51, 51);">跟：</font>

```java
list.forEach(new Consumer<String>() {
    @Override
    public void accept(String item) {
        System.out.println(item);
    }
});
```

<font style="color:rgb(51, 51, 51);">在性能上，</font>**几乎无差异**<font style="color:rgb(51, 51, 51);">，有时 lambda 甚至更快（因为 JVM 会做 Lambda 表达式的 </font><font style="color:rgb(51, 51, 51);">invokedynamic</font><font style="color:rgb(51, 51, 51);"> 优化）。</font>

---



<h3 id="d8e3a601">特殊场景下，Lambda 可能会“慢很多”</h3>

<font style="color:rgb(51, 51, 51);">说“Lambda 让 Java 慢 30 倍”的人，</font>**通常是指在一些特定场景下 Lambda 带来的性能坑**<font style="color:rgb(51, 51, 51);">，比如：</font>

<h4 id="0deb1349">场景一：Lambda 捕获了外部变量，导致频繁创建对象</h4>

```java
public List<Runnable> test() {
    List<Runnable> list = new ArrayList<>();
    for (int i = 0; i < 1000000; i++) {
        int finalI = i;
        list.add(() -> System.out.println(finalI));
    }
    return list;
}
```

<font style="color:#0e0e0e;"></font><font style="color:#0e0e0e;">每个 lambda 都捕获了不同的 i，所以 </font>**<font style="color:#0e0e0e;">每次都会创建一个新的函数对象（闭包）</font>**<font style="color:#0e0e0e;">，而不是重用。</font>

<font style="color:rgb(51, 51, 51);">相同逻辑的匿名类写法可能会复用对象，从而更节省内存。</font>

---

<h4 id="2fd07984">场景二：频繁在热点代码中使用 Lambda，会影响 JIT 优化</h4>

<font style="color:rgb(51, 51, 51);">一些性能敏感的循环内使用 Lambda，可能因为虚调用（</font><font style="color:rgb(51, 51, 51);">invokedynamic</font><font style="color:rgb(51, 51, 51);">）影响 JIT 编译器的内联优化，从而 </font>**无法达到预期的性能表现**<font style="color:rgb(51, 51, 51);">。</font>

---

<h4 id="575274d3">场景三：Stream + Lambda 的组合，某些时候不如手写循环快</h4>

```java
list.stream()
.filter(x -> x > 10)
.map(x -> x * 2)
.collect(Collectors.toList());
```

<font style="color:rgb(51, 51, 51);">这个 Stream 写法确实</font>**优雅**<font style="color:rgb(51, 51, 51);">，但在超大数据集（上百万条）上，</font>**比不上手写 for 循环那样零开销**<font style="color:rgb(51, 51, 51);">：</font>

```java
List<Integer> result = new ArrayList<>();
for (Integer x : list) {
    if (x > 10) result.add(x * 2);
}
```

<font style="color:rgb(51, 51, 51);">在微基准测试中，差距可以达到数倍甚至十几倍。</font>

---

<h3 id="1a525f3d">总结一句话：</h3>

**<font style="color:#0e0e0e;">Lambda 本身不是慢，而是在特定场景下会引入开销。如果你在意性能，就要了解背后的实现。</font>**

---

<h3 id="35fa54fc">建议</h3>

+ <font style="color:rgb(51, 51, 51);">写业务逻辑，优先 </font>**清晰 + 可读性**<font style="color:rgb(51, 51, 51);">，放心用 Lambda；</font>
+ <font style="color:rgb(51, 51, 51);">写性能关键代码（比如高频调用的函数），可用 JMH 做微基准测试；</font>
+ <font style="color:rgb(51, 51, 51);">遇到性能问题再优化，而不是一开始就避免 Lambda。</font>

---



<h1 id="rc7nt">四、Lambda 表达式总结</h1>
<h2 id="sgsoS">浅谈 Lambda 表达式</h2>

<font style="color:rgb(34, 34, 38);">lambda表达式可以写出更简洁的代码，之前在Java里面要传递一段逻辑如果没有Lambda表达式这种参数化代码的方式那就只能定义类和和创建对象或者使用匿名内部类来传递，代码上比较复杂。如果不使用Lambda表达式当然也可以实现同样的功能，但是在开发的过程中，一个项目肯定是多个人一起开发的，如果别人使用Lambda表达式，自己至少要能看懂。</font>  
<font style="color:rgb(34, 34, 38);">至于</font>`**<font style="color:rgb(34, 34, 38);">foreach</font>**`<font style="color:rgb(34, 34, 38);">效率问题，是Java8的另一个新特性Stream流功能，Stream流提供pallStream并行流可以在某些场景上提供并发操作的逻辑提升效率，但是也要注意线程安全问题。</font>

<font style="color:rgb(34, 34, 38);"></font>

<h2 id="G1te2"><font style="color:rgb(34, 34, 38);">什么时候更适合使用 Lambda 表达式</font></h2>

Lambda 表达式适合用于需要传递短小逻辑代码的场景，尤其是在函数式接口中，例如集合的遍历、排序、过滤等操作时，可以让代码更加简洁清晰，避免冗长的匿名内部类。它特别适合用于一次性、不需要复用的简单逻辑，比如线程的启动、事件监听、流式处理中的中间操作等。相比传统方式，Lambda 写法更直观，也更符合现代 Java 的风格。不过当逻辑较复杂、需要调试、或涉及性能敏感的循环操作时，传统写法可能更具可读性和可控性。因此，在日常开发中，是否使用 Lambda 表达式，应根据具体场景权衡简洁性、可读性和性能等因素，做到合理使用、适度使用。



<h2 id="EXy0G">结语</h2>

> 最后小王想跟大家说，Lambda 表达式的学习曲线是相对陡峭的，有的朋友也会跟我说，Lambda 表达式好优雅，有想去用，但是之前的方式写习惯了，动力不足，而且不知道，啥时候该用？
>
> 小王在这里发表一下自己的看法，很多老鸟在开发中，业务代码写 Lambda 表达式 + Stream 流，得心应手，但是像我们这种小白，比如像小王，在 Golang 转会 Java，本来就菜，这不更雪上加霜啊。。。所以我们也没有必要过多担心，首先要保证在开发中先看懂别人写的 Lambda 表达式，再去自己慢慢沉淀就好啦～

