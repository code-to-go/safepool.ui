<script>
  import {
    BlockTitle,
    f7,
    f7ready,
    Icon,
    Link,
    List,
    ListItem,
    Message,
    Messagebar,
    MessagebarAttachment,
    MessagebarAttachments,
    MessagebarSheet,
    MessagebarSheetImage,
    Messages,
    MessagesTitle,
    Navbar,
    NavLeft,
    NavRight,
    NavTitle,
    Page,
    theme,
  } from "framework7-svelte";
  import { onMount } from "svelte";
  import {
    GetPoolList,
    GetMessages,
    GetSelf,
    PostMessage,
  } from "../../wailsjs/go/main/App";
  import { pool } from "../../wailsjs/go/models";

  let messagebarInstance;

  export let poolName;

  let selfIdentity;
  let messages = [];
  let lastId = '9223372036854775807'
  onMount(async () => {
    selfIdentity = await GetSelf();
  })

  onMount(async () => {
    messages = (await GetMessages(poolName, '0', lastId, 32)) || [];
    console.log(messages)
  });

  let typingMessage = null;
  let messageText = "";

  onMount(() => {
    f7ready(() => {
      messagebarInstance = f7.messagebar.get(".messagebar");
    });
  });

  function isFirstMessage(message, index) {
    const previousMessage = messages[index - 1];
    if (message.isTitle) return false;
    if (
      !previousMessage ||
      previousMessage.type !== message.type ||
      previousMessage.name !== message.name
    )
      return true;
    return false;
  }
  function isLastMessage(message, index) {
    const nextMessage = messages[index + 1];
    if (message.isTitle) return false;
    if (
      !nextMessage ||
      nextMessage.type !== message.type ||
      nextMessage.name !== message.name
    )
      return true;
    return false;
  }
  function isTailMessage(message, index) {
    const nextMessage = messages[index + 1];
    if (message.isTitle) return false;
    if (
      !nextMessage ||
      nextMessage.type !== message.type ||
      nextMessage.name !== message.name
    )
      return true;
    return false;
  }
  async function sendMessage() {
    const text = messageText.replace(/\n/g, "<br>").trim();
    if (text.length) {
      const m = {
        author: selfIdentity.s.P,
        contentType: "text/html",
        content: text,
      }
      await PostMessage(poolName, m);
      console.log('sent message: ', m)
      messages.push(m)
    }
    // Clear
    messageText = "";
    messagebarInstance.clear();

    // Focus area
    if (text.length) messagebarInstance.focus();
  }
</script>

<Page class="page-chat">
  <Messagebar
    resizePage
    value={messageText}
    onInput={(e) => (messageText = e.target.value)}
  >
    <a class="link icon-only" slot="inner-end" on:click={sendMessage}>
      <Icon
        ios="f7:arrow_up_circle_fill"
        aurora="f7:arrow_up_circle_fill"
        md="material:send"
      />
    </a>
  </Messagebar>
  {#each messages as message}
    {message.content}
  {/each}

  <Messages>
    <MessagesTitle><b>Sunday, Feb 9,</b> 12:58</MessagesTitle>
    <Message
        type="received"
        htmlText="content"
      />
    {#each messages as message, index (index)}
    {message.content}
      <Message
        type={selfIdentity == message.author ? "received" : "sent"}
        name={selfIdentity == message.author ? null : message.author.nick}
        first={isFirstMessage(message, index)}
        last={isLastMessage(message, index)}
        tail={isTailMessage(message, index)}
        htmlText={message.content}
      />
    {/each}
    {#if typingMessage}
      <Message
        type="received"
        typing={true}
        first={true}
        last={true}
        tail={true}
        header={`${typingMessage.name} is typing`}
        avatar={typingMessage.avatar}
      />
    {/if}

  </Messages>
</Page>

<style>
</style>
