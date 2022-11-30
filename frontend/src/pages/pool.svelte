<script>
  import {
    Navbar,
    Page,
    Block,
    Tabs,
    Tab,
    Link,
    Toolbar,
    NavRight,
    Button,
  } from "framework7-svelte";
  import Chat from "./chat.svelte";
  import Library from "./library.svelte";

  export let poolName;
  let activeFeature = "chat";
</script>

<Page pageContent={false}>
  <Navbar title={poolName} backLink="Back">
    <NavRight>
      <Button active={activeFeature == 'chat'} onClick={()=> activeFeature = 'chat'}>Chat</Button>
      <Button active={activeFeature == "library"} onClick={()=> activeFeature = 'library'}>Library</Button>
      <Button active={activeFeature == "people"} onClick={()=> activeFeature = 'people'}>People</Button>
    </NavRight>
  </Navbar>
  {#if activeFeature == 'chat'}
    <Chat poolName={poolName}/>
  {/if}
  {#if activeFeature == 'library'}
    <Library poolName={poolName}/>
  {/if}
  {#if activeFeature == "people"}
    LATER
  {/if}
</Page>
