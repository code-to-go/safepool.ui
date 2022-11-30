<script>
  import {
    Page,
    Navbar,
    NavLeft,
    NavTitle,
    NavTitleLarge,
    NavRight,
    Link,
    Toolbar,
    Block,
    BlockTitle,
    List,
    ListItem,
    Row,
    Col,
    Button,
    Icon,
    Views,
    View,
    Tabs,
    Tab,
    theme,
    Panel,
    Popup,
  } from "framework7-svelte";
  import { onMount } from "svelte";

  import { GetPoolList } from "../../wailsjs/go/main/App";
  import Addpopup from "./addpopup.svelte";

  let pools = ["something not working!"];
  onMount(async () => {
    pools = await GetPoolList() || [];
  });

  let popupOpened = false;
</script>

<Addpopup {popupOpened} onPopupClosed={() => popupOpened = false}/>

<Page name="view-home">
  <!-- Views/Tabs container -->
  <Views tabs class="pool-areas">
    <!-- Tabbar for switching views-tabs -->
    <Toolbar tabbar labels bottom>
      <Link
        tabLink="#view-pool"
        tabLinkActive
        iconIos="f7:house_fill"
        iconAurora="f7:house_fill"
        iconMd="material:home"
        text="Home"
      />
      <Link
        tabLink="#view-settings"
        iconIos="f7:gear"
        iconAurora="f7:gear"
        iconMd="material:settings"
        text="Settings"
      />
    </Toolbar>

    <Tabs>
      <Tab id="view-pool" tabActive>
        <Navbar>
          <NavTitle>Safepool <Icon material="water" color="blue" /></NavTitle>
          <NavRight>
            <Button tabletFullscreen onClick={() => (popupOpened = true)}
              >Add</Button
            >
          </NavRight>
        </Navbar>
        <List class="components-list">
          {#each pools as pool}
            <ListItem
              reloadDetail={theme.aurora}
              link="/pools/{pool.replace('/', '|')}"
              title={pool}
            >
              <i class="icon icon-f7" slot="media" />
            </ListItem>
          {/each}
        </List>
      </Tab>
    </Tabs>

    <!-- Settings View -->
  </Views>
</Page>
