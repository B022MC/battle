declare namespace API {
  type ShopsGroupsBaseParams = {
    house_gid: number;
  };

  type ShopsGroupsForbidParams = ShopsGroupsBaseParams & {
    key: string;
    member_ids?: number[];
  };

  type ShopsGroupsBindParams = ShopsGroupsBaseParams & {
    message_id: number;
  };
}

