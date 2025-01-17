delete from account_operations where block_num>9090306;
delete from allowed_protocols where block_num>9090306;
delete from allowed_tokens where block_num>9090306;
delete from credit_manager_stats where block_num>9090306;
delete from credit_session_snapshots where block_num>9090306;
delete from dao_operations where block_num> 9090306;


delete from debts where block_num> 9090306;
update debt_sync set last_calculated_at=9090306;
delete from token_oracle where  block_num> 9090306;
delete from transfer_account_allowed  where  block_num> 9090306;
delete from no_session_transfers where  block_num> 9090306;
delete from fast_check_params where  block_num> 9090306;
delete from parameters where  block_num> 9090306;
delete from pool_stats where  block_num> 9090306;
delete from pool_ledger where  block_num> 9090306;
delete from current_debts where block_num > 9090306;  
delete from liquidable_accounts where block_num> 9090306;

-- a
update sync_adapters set last_sync = 9090306 where type in ('CreditConfigurator',  'AccountManager','CreditFilter', 'CreditManager', 'Pool') and last_sync > 9090306;
-- 
delete from credit_sessions where since> 9090306;
update  credit_sessions set closed_at=0,liquidator='', remaining_funds='0', status=0 where closed_at> 9090306;

-- the borrower is not accounted for ; i.e. transfer of acocunt is not handled and should be used carefully
update credit_sessions set borrowed_amount=css.borrowed_amount_bi, balances=css.balances from  credit_sessions cs join (select distinct on(session_id) * from credit_session_snapshots order by session_id, block_num desc) css on css.session_id=cs.id  where cs.status=0;

-- for creditmanager

update credit_managers cm set 
    opened_accounts_count=cms.opened_accounts_count,
    total_opened_accounts=cms.total_opened_accounts,
    total_repaid_accounts=cms.total_repaid_accounts,
    total_closed_accounts=cms.total_closed_accounts,
    total_liquidated_accounts=cms.total_liquidated_accounts,
    total_borrowed_bi=cms.total_borrowed_bi,
    cumulative_borrowed_bi=cms.cumulative_borrowed_bi,
    available_liquidity_bi=cms.available_liquidity_bi,
    total_repaid_bi=cms.total_repaid_bi,
    total_losses_bi=cms.total_losses_bi,
    total_profit_bi=cms.total_profit_bi,
    total_borrowed=cms.total_borrowed,
    cumulative_borrowed=cms.cumulative_borrowed,
    available_liquidity=cms.available_liquidity,
    total_repaid=cms.total_repaid,
    total_losses=cms.total_losses,
    total_profit=cms.total_profit,
    unique_users=cms.unique_users from (select distinct on (credit_manager) * from credit_manager_stats order by credit_manager, block_num desc) cms 
    where cm.address= cms.credit_manager;




----
-- delete from sync_adapters where discovered_at> 9090306;
-- delete from price_feeds where block_num >  9090306;
-- delete from token_oracle where block_num > 9090306;
-- update sync_adapters set last_sync=9090306 where type in ('PriceOracle', 'ChainlinkPriceFeed', 'CompositeChainlinkPF', 'QueryPriceFeed');